# UpdateEmployeeData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The UUID of the employee to be updated. | 
**Type** | **string** |  | 
**Attributes** | [**UpdateEmployeeDataAttributes**](UpdateEmployeeDataAttributes.md) |  | 

## Methods

### NewUpdateEmployeeData

`func NewUpdateEmployeeData(id string, type_ string, attributes UpdateEmployeeDataAttributes, ) *UpdateEmployeeData`

NewUpdateEmployeeData instantiates a new UpdateEmployeeData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateEmployeeDataWithDefaults

`func NewUpdateEmployeeDataWithDefaults() *UpdateEmployeeData`

NewUpdateEmployeeDataWithDefaults instantiates a new UpdateEmployeeData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateEmployeeData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateEmployeeData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateEmployeeData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateEmployeeData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateEmployeeData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateEmployeeData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateEmployeeData) GetAttributes() UpdateEmployeeDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateEmployeeData) GetAttributesOk() (*UpdateEmployeeDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateEmployeeData) SetAttributes(v UpdateEmployeeDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


