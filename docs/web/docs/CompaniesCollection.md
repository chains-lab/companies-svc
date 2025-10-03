# CompaniesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]CompanyData**](CompanyData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewCompaniesCollection

`func NewCompaniesCollection(data []CompanyData, links PaginationData, ) *CompaniesCollection`

NewCompaniesCollection instantiates a new CompaniesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompaniesCollectionWithDefaults

`func NewCompaniesCollectionWithDefaults() *CompaniesCollection`

NewCompaniesCollectionWithDefaults instantiates a new CompaniesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CompaniesCollection) GetData() []CompanyData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CompaniesCollection) GetDataOk() (*[]CompanyData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CompaniesCollection) SetData(v []CompanyData)`

SetData sets Data field to given value.


### GetLinks

`func (o *CompaniesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *CompaniesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *CompaniesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


