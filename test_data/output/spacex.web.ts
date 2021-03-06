/* tslint:disable */
/* eslint-disable */
// @ts-ignore
import { WebDefinitions, OperationDefinition, Definitions as AllDefinitions } from '@openland/spacex';
// @ts-ignore
const { list, notNull, scalar, field, obj, inline, fragment, args, fieldValue, refValue, intValue, floatValue, stringValue, boolValue, listValue, objectValue } = WebDefinitions;

const DeviceNanoSelector = obj(
            field('__typename', '__typename', args(), notNull(scalar('String'))),
            field('id', 'id', args(), notNull(scalar('ID')))
        );

const LampShortSelector = obj(
            field('__typename', '__typename', args(), notNull(scalar('String'))),
            field('id', 'id', args(), notNull(scalar('ID'))),
            field('brightness', 'brightness', args(), notNull(scalar('Float'))),
            field('minBrightness', 'minBrightness', args(), notNull(scalar('Float'))),
            field('maxBrightness', 'maxBrightness', args(), notNull(scalar('Float')))
        );

const DeviceShortSelector = obj(
            field('__typename', '__typename', args(), notNull(scalar('String'))),
            field('description', 'description', args(), scalar('String')),
            field('addedBy', 'addedBy', args(), notNull(obj(
                    field('__typename', '__typename', args(), notNull(scalar('String'))),
                    field('id', 'id', args(), notNull(scalar('ID'))),
                    field('username', 'username', args(), notNull(scalar('String')))
                ))),
            fragment('Device', DeviceNanoSelector),
            inline('Lamp', obj(
                field('__typename', '__typename', args(), notNull(scalar('String'))),
                fragment('Lamp', LampShortSelector)
            )),
            inline('Switch', obj(
                field('__typename', '__typename', args(), notNull(scalar('String'))),
                field('id', 'id', args(), notNull(scalar('ID'))),
                field('numberOfButtons', 'numberOfButtons', args(), scalar('Int'))
            )),
            inline('Lamp', obj(
                field('__typename', '__typename', args(), notNull(scalar('String'))),
                field('id', 'id', args(), notNull(scalar('ID'))),
                field('brightness', 'brightness', args(), notNull(scalar('Float')))
            ))
        );

const UserShortSelector = obj(
            field('__typename', '__typename', args(), notNull(scalar('String'))),
            field('id', 'id', args(), notNull(scalar('ID')))
        );

const DiscoverDevicesSelector = obj(
            field('discover', 'discover', args(), notNull(list(notNull(obj(
                    field('__typename', '__typename', args(), notNull(scalar('String'))),
                    field('id', 'id', args(), notNull(scalar('ID'))),
                    field('name', 'name', args(), notNull(scalar('String'))),
                    field('host', 'host', args(), notNull(scalar('String'))),
                    field('port', 'port', args(), notNull(scalar('Int')))
                )))))
        );
const AddUserSelector = obj(
            field('addUser', 'addUser', args(fieldValue("username", refValue('username')), fieldValue("password", refValue('password'))), notNull(obj(
                    field('__typename', '__typename', args(), notNull(scalar('String'))),
                    fragment('User', UserShortSelector)
                )))
        );
const UpdateDeviceDescriptionSelector = obj(
            field('updateDeviceDescription', 'updateDeviceDescription', args(fieldValue("id", refValue('id')), fieldValue("description", refValue('description'))), notNull(obj(
                    field('__typename', '__typename', args(), notNull(scalar('String'))),
                    fragment('Device', DeviceShortSelector)
                )))
        );
export const Operations: { [key: string]: OperationDefinition } = {
    DiscoverDevices: {
        kind: 'query',
        name: 'DiscoverDevices',
        body: 'query DiscoverDevices{discover{__typename id name host port}}',
        selector: DiscoverDevicesSelector
    },
    AddUser: {
        kind: 'mutation',
        name: 'AddUser',
        body: 'mutation AddUser($username:String!,$password:String!){addUser(username:$username,password:$password){__typename ...UserShort}}fragment UserShort on User{__typename id}',
        selector: AddUserSelector
    },
    UpdateDeviceDescription: {
        kind: 'mutation',
        name: 'UpdateDeviceDescription',
        body: 'mutation UpdateDeviceDescription($id:ID!,$description:DeviceDescription!){updateDeviceDescription(id:$id,description:$description){__typename ...DeviceShort}}fragment DeviceShort on Device{__typename description addedBy{__typename id username}...DeviceNano ... on Lamp{__typename ...LampShort}... on Switch{__typename id numberOfButtons}... on Lamp{__typename id brightness}}fragment DeviceNano on Device{__typename id}fragment LampShort on Lamp{__typename id brightness minBrightness maxBrightness}',
        selector: UpdateDeviceDescriptionSelector
    },
};
export const Definitions: AllDefinitions = { operations: Operations };