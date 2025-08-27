# EmployeesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]EmployeeData**](EmployeeData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewEmployeesCollection

`func NewEmployeesCollection(data []EmployeeData, links PaginationData, ) *EmployeesCollection`

NewEmployeesCollection instantiates a new EmployeesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmployeesCollectionWithDefaults

`func NewEmployeesCollectionWithDefaults() *EmployeesCollection`

NewEmployeesCollectionWithDefaults instantiates a new EmployeesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *EmployeesCollection) GetData() []EmployeeData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *EmployeesCollection) GetDataOk() (*[]EmployeeData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *EmployeesCollection) SetData(v []EmployeeData)`

SetData sets Data field to given value.


### GetLinks

`func (o *EmployeesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *EmployeesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *EmployeesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


