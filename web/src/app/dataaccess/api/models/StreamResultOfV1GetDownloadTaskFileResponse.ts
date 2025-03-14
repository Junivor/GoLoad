/* tslint:disable */
/* eslint-disable */
/**
 * go_load/v1/go_load.proto
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: version not set
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from '../runtime';
import type { V1GetDownloadTaskFileResponse } from './V1GetDownloadTaskFileResponse';
import {
    V1GetDownloadTaskFileResponseFromJSON,
    V1GetDownloadTaskFileResponseFromJSONTyped,
    V1GetDownloadTaskFileResponseToJSON,
    V1GetDownloadTaskFileResponseToJSONTyped,
} from './V1GetDownloadTaskFileResponse';
import type { RpcStatus } from './RpcStatus';
import {
    RpcStatusFromJSON,
    RpcStatusFromJSONTyped,
    RpcStatusToJSON,
    RpcStatusToJSONTyped,
} from './RpcStatus';

/**
 * 
 * @export
 * @interface StreamResultOfV1GetDownloadTaskFileResponse
 */
export interface StreamResultOfV1GetDownloadTaskFileResponse {
    /**
     * 
     * @type {V1GetDownloadTaskFileResponse}
     * @memberof StreamResultOfV1GetDownloadTaskFileResponse
     */
    result?: V1GetDownloadTaskFileResponse;
    /**
     * 
     * @type {RpcStatus}
     * @memberof StreamResultOfV1GetDownloadTaskFileResponse
     */
    error?: RpcStatus;
}

/**
 * Check if a given object implements the StreamResultOfV1GetDownloadTaskFileResponse interface.
 */
export function instanceOfStreamResultOfV1GetDownloadTaskFileResponse(value: object): value is StreamResultOfV1GetDownloadTaskFileResponse {
    return true;
}

export function StreamResultOfV1GetDownloadTaskFileResponseFromJSON(json: any): StreamResultOfV1GetDownloadTaskFileResponse {
    return StreamResultOfV1GetDownloadTaskFileResponseFromJSONTyped(json, false);
}

export function StreamResultOfV1GetDownloadTaskFileResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): StreamResultOfV1GetDownloadTaskFileResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'result': json['result'] == null ? undefined : V1GetDownloadTaskFileResponseFromJSON(json['result']),
        'error': json['error'] == null ? undefined : RpcStatusFromJSON(json['error']),
    };
}

export function StreamResultOfV1GetDownloadTaskFileResponseToJSON(json: any): StreamResultOfV1GetDownloadTaskFileResponse {
    return StreamResultOfV1GetDownloadTaskFileResponseToJSONTyped(json, false);
}

export function StreamResultOfV1GetDownloadTaskFileResponseToJSONTyped(value?: StreamResultOfV1GetDownloadTaskFileResponse | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'result': V1GetDownloadTaskFileResponseToJSON(value['result']),
        'error': RpcStatusToJSON(value['error']),
    };
}

