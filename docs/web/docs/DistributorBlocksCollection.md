# DistributorBlocksCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]DistributorBlockData**](DistributorBlockData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewDistributorBlocksCollection

`func NewDistributorBlocksCollection(data []DistributorBlockData, links PaginationData, ) *DistributorBlocksCollection`

NewDistributorBlocksCollection instantiates a new DistributorBlocksCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDistributorBlocksCollectionWithDefaults

`func NewDistributorBlocksCollectionWithDefaults() *DistributorBlocksCollection`

NewDistributorBlocksCollectionWithDefaults instantiates a new DistributorBlocksCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *DistributorBlocksCollection) GetData() []DistributorBlockData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DistributorBlocksCollection) GetDataOk() (*[]DistributorBlockData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DistributorBlocksCollection) SetData(v []DistributorBlockData)`

SetData sets Data field to given value.


### GetLinks

`func (o *DistributorBlocksCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *DistributorBlocksCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *DistributorBlocksCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


