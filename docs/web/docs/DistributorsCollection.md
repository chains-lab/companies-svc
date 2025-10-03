# companiesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]companyData**](companyData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewcompaniesCollection

`func NewcompaniesCollection(data []companyData, links PaginationData, ) *companiesCollection`

NewcompaniesCollection instantiates a new companiesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewcompaniesCollectionWithDefaults

`func NewcompaniesCollectionWithDefaults() *companiesCollection`

NewcompaniesCollectionWithDefaults instantiates a new companiesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *companiesCollection) GetData() []companyData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *companiesCollection) GetDataOk() (*[]companyData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *companiesCollection) SetData(v []companyData)`

SetData sets Data field to given value.


### GetLinks

`func (o *companiesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *companiesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *companiesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


