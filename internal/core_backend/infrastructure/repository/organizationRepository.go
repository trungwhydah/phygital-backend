package repository

import (
	"context"
	"errors"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// OrganizationRepository struct
type OrganizationRepository struct {
	dbMongo *mongo.Database
}

// NewOrganizationRepository create repository
func NewOrganizationRepository(dbMongo *mongo.Database) *OrganizationRepository {
	return &OrganizationRepository{dbMongo: dbMongo}
}

func (r *OrganizationRepository) CreateOrganization(organization *entity.Organization) (*entity.Organization, error) {
	result, err := r.dbMongo.Collection(organization.CollectionName()).InsertOne(context.TODO(), &organization)
	if err != nil {
		return nil, err
	}
	organization.ID = result.InsertedID.(primitive.ObjectID)

	return organization, nil
}

func (r *OrganizationRepository) GetAllOrganizations() (*[]entity.Organization, error) {

	cursor, err := r.dbMongo.Collection(entity.Organization{}.CollectionName()).Find(context.TODO(), bson.M{"status": common.StatusActive})
	if err != nil {
		return nil, err
	}

	var orgs []entity.Organization
	if err = cursor.All(context.TODO(), &orgs); err != nil {
		return nil, err
	}

	return &orgs, nil
}

// UpdateOrganization
func (r *OrganizationRepository) UpdateOrganization(organization *entity.Organization) (bool, error) {
	filter := bson.D{{Key: "_id", Value: organization.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "org_name", Value: organization.OrganizationName},
			{Key: "org_tag_name", Value: organization.NameTag},
			{Key: "org_logo_url", Value: organization.LogoURL},
		}}}
	result, err := r.dbMongo.Collection(organization.CollectionName()).UpdateOne(
		context.TODO(),
		filter,
		update)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount+result.UpsertedCount+result.MatchedCount != 0, nil
}

// GetOrgDetail
func (r *OrganizationRepository) GetDetailOrganization(orgID *string) (*entity.Organization, error) {
	oID, err := primitive.ObjectIDFromHex(*orgID)
	if err != nil {
		return nil, err
	}

	var org entity.Organization
	err = r.dbMongo.Collection(org.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: oID}}).Decode(&org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// GetOrgByTagName
func (r *OrganizationRepository) GetOrgByTagName(tagName *string) (*entity.Organization, error) {
	var org entity.Organization
	err := r.dbMongo.Collection(org.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "org_tag_name", Value: *tagName}}).Decode(&org)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &org, nil
}
