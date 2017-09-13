import { Injectable } from '@angular/core';

import { Image } from '../entities/image';

import { ApiService } from './api.service';

@Injectable()
export class ImageService {
    constructor(private apiService: ApiService) { }

    getImages(): Promise<Image[]> {
        return this.apiService.get('/images/json')
            .then(response => response.json() as Image[])
            .then(images => {
                return images.map(image => {
                    if (image.RepoTags === null) {
                        [image.Repository] = image.RepoDigests[0].split('@')
                        image.Tag = '<none>';
                    } else {
                        [image.Repository, image.Tag] = image.RepoTags[0].split(':');
                    }

                    return image;
                })
            });
    }

    prune(): Promise<any> {
        return this.apiService.post('/images/prune?filters={"dangling": true}', null)
            .then(response => response.json())
    }

    pull(name: string): Promise<any> {
        return this.apiService.post(`/images/create?fromImage=${name}`, null);
    }

    remove(id: string): Promise<any> {
        return this.apiService.delete(`/images/${id}`);
    }
}
