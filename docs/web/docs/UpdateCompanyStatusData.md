# UpdateCompanyStatusData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | The UUID of the company to be updated. | 
**Type** | **string** |  | 
**Attributes** | [**UpdateCompanyStatusDataAttributes**](UpdateCompanyStatusDataAttributes.md) |  | 

## Methods

### NewUpdateCompanyStatusData

`func NewUpdateCompanyStatusData(id uuid.UUID, type_ string, attributes UpdateCompanyStatusDataAttributes, ) *UpdateCompanyStatusData`

NewUpdateCompanyStatusData instantiates a new UpdateCompanyStatusData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateCompanyStatusDataWithDefaults

`func NewUpdateCompanyStatusDataWithDefaults() *UpdateCompanyStatusData`

NewUpdateCompanyStatusDataWithDefaults instantiates a new UpdateCompanyStatusData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateCompanyStatusData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateCompanyStatusData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateCompanyStatusData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateCompanyStatusData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateCompanyStatusData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateCompanyStatusData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateCompanyStatusData) GetAttributes() UpdateCompanyStatusDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateCompanyStatusData) GetAttributesOk() (*UpdateCompanyStatusDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateCompanyStatusData) SetAttributes(v UpdateCompanyStatusDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


