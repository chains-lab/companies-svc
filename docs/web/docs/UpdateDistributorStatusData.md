# UpdateDistributorStatusData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The UUID of the distributor to be updated. | 
**Type** | **string** |  | 
**Attributes** | [**UpdateDistributorStatusDataAttributes**](UpdateDistributorStatusDataAttributes.md) |  | 

## Methods

### NewUpdateDistributorStatusData

`func NewUpdateDistributorStatusData(id string, type_ string, attributes UpdateDistributorStatusDataAttributes, ) *UpdateDistributorStatusData`

NewUpdateDistributorStatusData instantiates a new UpdateDistributorStatusData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateDistributorStatusDataWithDefaults

`func NewUpdateDistributorStatusDataWithDefaults() *UpdateDistributorStatusData`

NewUpdateDistributorStatusDataWithDefaults instantiates a new UpdateDistributorStatusData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateDistributorStatusData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateDistributorStatusData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateDistributorStatusData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateDistributorStatusData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateDistributorStatusData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateDistributorStatusData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateDistributorStatusData) GetAttributes() UpdateDistributorStatusDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateDistributorStatusData) GetAttributesOk() (*UpdateDistributorStatusDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateDistributorStatusData) SetAttributes(v UpdateDistributorStatusDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


