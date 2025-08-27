# EmployeeInviteData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | employee invite id | 
**Type** | **string** |  | 
**Attributes** | [**EmployeeInviteAttributes**](EmployeeInviteAttributes.md) |  | 

## Methods

### NewEmployeeInviteData

`func NewEmployeeInviteData(id string, type_ string, attributes EmployeeInviteAttributes, ) *EmployeeInviteData`

NewEmployeeInviteData instantiates a new EmployeeInviteData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmployeeInviteDataWithDefaults

`func NewEmployeeInviteDataWithDefaults() *EmployeeInviteData`

NewEmployeeInviteDataWithDefaults instantiates a new EmployeeInviteData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *EmployeeInviteData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *EmployeeInviteData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *EmployeeInviteData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *EmployeeInviteData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *EmployeeInviteData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *EmployeeInviteData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *EmployeeInviteData) GetAttributes() EmployeeInviteAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *EmployeeInviteData) GetAttributesOk() (*EmployeeInviteAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *EmployeeInviteData) SetAttributes(v EmployeeInviteAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


