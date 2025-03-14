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
import type { V1DownloadTask } from './V1DownloadTask';
import {
    V1DownloadTaskFromJSON,
    V1DownloadTaskFromJSONTyped,
    V1DownloadTaskToJSON,
    V1DownloadTaskToJSONTyped,
} from './V1DownloadTask';

/**
 * 
 * @export
 * @interface V1GetDownloadTaskListResponse
 */
export interface V1GetDownloadTaskListResponse {
    /**
     * 
     * @type {Array<V1DownloadTask>}
     * @memberof V1GetDownloadTaskListResponse
     */
    downloadTaskList?: Array<V1DownloadTask>;
    /**
     * 
     * @type {string}
     * @memberof V1GetDownloadTaskListResponse
     */
    totalDownloadTaskCount?: string;
}

/**
 * Check if a given object implements the V1GetDownloadTaskListResponse interface.
 */
export function instanceOfV1GetDownloadTaskListResponse(value: object): value is V1GetDownloadTaskListResponse {
    return true;
}

export function V1GetDownloadTaskListResponseFromJSON(json: any): V1GetDownloadTaskListResponse {
    return V1GetDownloadTaskListResponseFromJSONTyped(json, false);
}

export function V1GetDownloadTaskListResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1GetDownloadTaskListResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'downloadTaskList': json['downloadTaskList'] == null ? undefined : ((json['downloadTaskList'] as Array<any>).map(V1DownloadTaskFromJSON)),
        'totalDownloadTaskCount': json['totalDownloadTaskCount'] == null ? undefined : json['totalDownloadTaskCount'],
    };
}

export function V1GetDownloadTaskListResponseToJSON(json: any): V1GetDownloadTaskListResponse {
    return V1GetDownloadTaskListResponseToJSONTyped(json, false);
}

export function V1GetDownloadTaskListResponseToJSONTyped(value?: V1GetDownloadTaskListResponse | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'downloadTaskList': value['downloadTaskList'] == null ? undefined : ((value['downloadTaskList'] as Array<any>).map(V1DownloadTaskToJSON)),
        'totalDownloadTaskCount': value['totalDownloadTaskCount'],
    };
}

