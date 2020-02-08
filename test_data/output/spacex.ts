/* tslint:disable */
/* eslint-disable */
import * as Types from './spacex.types';
import { GraphqlEngine, GraphqlActiveSubscription, OperationParameters, GraphqlSubscriptionHandler, BaseSpaceXClient, SpaceQueryWatchParameters } from '@openland/spacex';

export class  extends BaseSpaceXClient {
    constructor(engine: GraphqlEngine) {
        super(engine);
    }
    queryDiscoverDevices(opts?: OperationParameters): Promise<Types.DiscoverDevices> {
        return this.query('DiscoverDevices', undefined, opts);
    }
    refetchDiscoverDevices(opts?: OperationParameters): Promise<Types.DiscoverDevices> {
        return this.refetch('DiscoverDevices', undefined);
    }
    updateDiscoverDevices(updater: (data: Types.DiscoverDevices) => Types.DiscoverDevices | null): Promise<boolean> {
        return this.updateQuery(updater, 'DiscoverDevices', undefined);
    }
    useDiscoverDevices(opts: SpaceQueryWatchParameters & { suspense: false }): Types.DiscoverDevices | null;
    useDiscoverDevices(opts: SpaceQueryWatchParameters): Types.DiscoverDevices;
    useDiscoverDevices(opts: SpaceQueryWatchParameters): Types.DiscoverDevices | null {
        return this.useQuery('DiscoverDevices', undefined, opts);
    }
    mutateAddUser(variables: Types.AddUserVariables): Promise<Types.AddUser> {
        return this.mutate('AddUser', variables);
    }
    mutateUpdateDeviceDescription(variables: Types.UpdateDeviceDescriptionVariables): Promise<Types.UpdateDeviceDescription> {
        return this.mutate('UpdateDeviceDescription', variables);
    }
}