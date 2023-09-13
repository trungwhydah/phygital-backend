// SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.9;

import "openzeppelin/contracts/token/ERC721/ERC721.sol";
import "openzeppelin/contracts/access/Ownable.sol";
import "openzeppelin/contracts/utils/Counters.sol";
import "openzeppelin/contracts/utils/Strings.sol";

contract DaNonNuocNFT is ERC721, Ownable {    
    using Strings for uint256;
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;
    string private baseURL;
    
    string public contractURIString;

    constructor(string memory _baseURL) ERC721("Da Non Nuoc", "DNN") {
        baseURL = _baseURL;
    }

    function safeMint(address to) public onlyOwner {
        require(to != address(0), "Cannot mint to zero address");
        
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        
        _safeMint(to, tokenId);
    }

    function _baseURI() internal view override returns (string memory) {
        return baseURL;
    }

    function setBaseURI(string memory _uri) external onlyOwner {
        baseURL = _uri;
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory)
    {
        
        require(_exists(tokenId), "ERC721URIStorage: URI query for nonexistent token");
  
        return string(abi.encodePacked(_baseURI(), tokenId.toString()));
    }
}
