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
/**
 * 
 * @export
 * @interface V1DeleteDownloadTaskRequest
 */
export interface V1DeleteDownloadTaskRequest {
    /**
     * 
     * @type {string}
     * @memberof V1DeleteDownloadTaskRequest
     */
    downloadTaskId?: string;
}

/**
 * Check if a given object implements the V1DeleteDownloadTaskRequest interface.
 */
export function instanceOfV1DeleteDownloadTaskRequest(value: object): value is V1DeleteDownloadTaskRequest {
    return true;
}

export function V1DeleteDownloadTaskRequestFromJSON(json: any): V1DeleteDownloadTaskRequest {
    return V1DeleteDownloadTaskRequestFromJSONTyped(json, false);
}

export function V1DeleteDownloadTaskRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1DeleteDownloadTaskRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'downloadTaskId': json['downloadTaskId'] == null ? undefined : json['downloadTaskId'],
    };
}

export function V1DeleteDownloadTaskRequestToJSON(json: any): V1DeleteDownloadTaskRequest {
    return V1DeleteDownloadTaskRequestToJSONTyped(json, false);
}

export function V1DeleteDownloadTaskRequestToJSONTyped(value?: V1DeleteDownloadTaskRequest | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'downloadTaskId': value['downloadTaskId'],
    };
}

