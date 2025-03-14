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
 * @interface V1UpdateDownloadTaskResponse
 */
export interface V1UpdateDownloadTaskResponse {
    /**
     * 
     * @type {V1DownloadTask}
     * @memberof V1UpdateDownloadTaskResponse
     */
    downloadTask?: V1DownloadTask;
}

/**
 * Check if a given object implements the V1UpdateDownloadTaskResponse interface.
 */
export function instanceOfV1UpdateDownloadTaskResponse(value: object): value is V1UpdateDownloadTaskResponse {
    return true;
}

export function V1UpdateDownloadTaskResponseFromJSON(json: any): V1UpdateDownloadTaskResponse {
    return V1UpdateDownloadTaskResponseFromJSONTyped(json, false);
}

export function V1UpdateDownloadTaskResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1UpdateDownloadTaskResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'downloadTask': json['downloadTask'] == null ? undefined : V1DownloadTaskFromJSON(json['downloadTask']),
    };
}

export function V1UpdateDownloadTaskResponseToJSON(json: any): V1UpdateDownloadTaskResponse {
    return V1UpdateDownloadTaskResponseToJSONTyped(json, false);
}

export function V1UpdateDownloadTaskResponseToJSONTyped(value?: V1UpdateDownloadTaskResponse | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'downloadTask': V1DownloadTaskToJSON(value['downloadTask']),
    };
}

