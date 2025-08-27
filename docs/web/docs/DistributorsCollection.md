# DistributorsCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]DistributorData**](DistributorData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewDistributorsCollection

`func NewDistributorsCollection(data []DistributorData, links PaginationData, ) *DistributorsCollection`

NewDistributorsCollection instantiates a new DistributorsCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDistributorsCollectionWithDefaults

`func NewDistributorsCollectionWithDefaults() *DistributorsCollection`

NewDistributorsCollectionWithDefaults instantiates a new DistributorsCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *DistributorsCollection) GetData() []DistributorData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DistributorsCollection) GetDataOk() (*[]DistributorData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DistributorsCollection) SetData(v []DistributorData)`

SetData sets Data field to given value.


### GetLinks

`func (o *DistributorsCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *DistributorsCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *DistributorsCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


