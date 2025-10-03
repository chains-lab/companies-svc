# CompanyBlocksCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]CompanyBlockData**](CompanyBlockData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewCompanyBlocksCollection

`func NewCompanyBlocksCollection(data []CompanyBlockData, links PaginationData, ) *CompanyBlocksCollection`

NewCompanyBlocksCollection instantiates a new CompanyBlocksCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompanyBlocksCollectionWithDefaults

`func NewCompanyBlocksCollectionWithDefaults() *CompanyBlocksCollection`

NewCompanyBlocksCollectionWithDefaults instantiates a new CompanyBlocksCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CompanyBlocksCollection) GetData() []CompanyBlockData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CompanyBlocksCollection) GetDataOk() (*[]CompanyBlockData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CompanyBlocksCollection) SetData(v []CompanyBlockData)`

SetData sets Data field to given value.


### GetLinks

`func (o *CompanyBlocksCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *CompanyBlocksCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *CompanyBlocksCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


