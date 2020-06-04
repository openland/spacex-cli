/* tslint:disable */
/* eslint-disable */
import * as Types from './spacex.types';
import { SpaceXClientParameters, GraphqlActiveSubscription, QueryParameters, MutationParameters, SubscriptionParameters, GraphqlSubscriptionHandler, BaseSpaceXClient, SpaceQueryWatchParameters } from '@openland/spacex';

export class ApiClient extends BaseSpaceXClient {
    constructor(params: SpaceXClientParameters) {
        super(params);
    }
    withParameters(params: Partial<SpaceXClientParameters>) {
        return new ApiClient({ ... params, engine: this.engine, globalCache: this.globalCache});
    }
    queryDiscoverDevices(params?: QueryParameters): Promise<Types.DiscoverDevices> {
        return this.query('DiscoverDevices', undefined, params);
    }
    refetchDiscoverDevices(params?: QueryParameters): Promise<Types.DiscoverDevices> {
        return this.refetch('DiscoverDevices', undefined, params);
    }
    updateDiscoverDevices(updater: (data: Types.DiscoverDevices) => Types.DiscoverDevices | null): Promise<boolean> {
        return this.updateQuery(updater, 'DiscoverDevices', undefined);
    }
    useDiscoverDevices(params: SpaceQueryWatchParameters & { suspense: false }): Types.DiscoverDevices | null;
    useDiscoverDevices(params?: SpaceQueryWatchParameters): Types.DiscoverDevices;
    useDiscoverDevices(params?: SpaceQueryWatchParameters): Types.DiscoverDevices | null {
        return this.useQuery('DiscoverDevices', undefined, params);
    }
    mutateAddUser(variables: Types.AddUserVariables, params?: MutationParameters): Promise<Types.AddUser> {
        return this.mutate('AddUser', variables, params);
    }
    mutateUpdateDeviceDescription(variables: Types.UpdateDeviceDescriptionVariables, params?: MutationParameters): Promise<Types.UpdateDeviceDescription> {
        return this.mutate('UpdateDeviceDescription', variables, params);
    }
}