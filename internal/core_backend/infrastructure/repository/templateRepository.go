package repository

import (
	"context"
	"errors"

	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateRepository struct {
	dbMongo *mongo.Database
}

func NewTemplateRepository(dbMongo *mongo.Database) *TemplateRepository {
	return &TemplateRepository{
		dbMongo: dbMongo,
	}
}

func (r *TemplateRepository) CreateTemplate(template *entity.Template) (*entity.Template, error) {
	result, err := r.dbMongo.Collection(template.CollectionName()).InsertOne(context.TODO(), &template)
	if err != nil {
		return nil, err
	}
	template.ID = result.InsertedID.(primitive.ObjectID)
	return template, nil
}

func (r *TemplateRepository) UpdateTemplate(template *entity.Template) error {
	update := bson.M{
		"$set": template,
	}

	filter := bson.M{
		"_id": template.ID,
	}

	_, err := r.dbMongo.Collection(template.CollectionName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *TemplateRepository) GetTemplate(ID *string) (*entity.Template, error) {
	templateID, err := primitive.ObjectIDFromHex(*ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: templateID}}
	var template entity.Template
	err = r.dbMongo.Collection(template.CollectionName()).FindOne(context.TODO(), filter).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) GetTemplateWebpages(tID *string) (*entity.TemplateWebpages, error) {
	templateID, err := primitive.ObjectIDFromHex(*tID)
	if err != nil {
		return nil, err
	}
	webpageColl := entity.WebPage{}.CollectionName()
	// Open an aggregation cursor
	coll := r.dbMongo.Collection(entity.Template{}.CollectionName())
	cursor, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{{"$match", bson.D{{"_id", templateID}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", webpageColl},
					{"localField", "pages.page_id"},
					{"foreignField", "_id"},
					{"as", "doc_pages"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", webpageColl},
					{"localField", "menu.page_id"},
					{"foreignField", "_id"},
					{"as", "doc_menu"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"pages",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$pages"},
									{"as", "pages"},
									{"in",
										bson.D{
											{"$first",
												bson.D{
													{"$filter",
														bson.D{
															{"input", "$doc_pages"},
															{"cond",
																bson.D{
																	{"$eq",
																		bson.A{
																			"$$this._id",
																			"$$pages.page_id",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{"menu",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$menu"},
									{"as", "menu"},
									{"in",
										bson.D{
											{"$mergeObjects",
												bson.A{
													bson.D{{"title", "$$menu.title"}},
													bson.D{
														{"$first",
															bson.D{
																{"$filter",
																	bson.D{
																		{"input", "$doc_menu"},
																		{"cond",
																			bson.D{
																				{"$eq",
																					bson.A{
																						"$$this._id",
																						"$$menu.page_id",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$unset",
				bson.A{
					"doc_pages",
					"doc_menu",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var templateWebpages []entity.TemplateWebpages
	if err = cursor.All(context.TODO(), &templateWebpages); err != nil {
		return nil, err
	}
	if len(templateWebpages) != 1 {
		err = errors.New("There must be exactly 1 template webpages found! But there are != 1!")
		return nil, err
	}
	return &templateWebpages[0], nil
}

func (r *TemplateRepository) CheckExistedTemplate(templateID *string) (bool, error) {
	templateObjectID, err := primitive.ObjectIDFromHex(*templateID)
	if err != nil {
		return false, err
	}
	filter := bson.D{{Key: "_id", Value: templateObjectID}}
	count, err := r.dbMongo.Collection(entity.Template{}.CollectionName()).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (r *TemplateRepository) GetAllTemplates() (*[]entity.Template, error) {
	filter := bson.M{} // Empty filter to retrieve all documents

	cursor, err := r.dbMongo.Collection(entity.Template{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	templates := []entity.Template{}
	for cursor.Next(context.TODO()) {
		var template entity.Template
		if err := cursor.Decode(&template); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &templates, nil
}
