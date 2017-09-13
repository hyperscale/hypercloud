import { Component, OnInit } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { Image } from '../../entities/image';

import { ImageService } from '../../services/image.service';
import { EventService } from '../../services/event.service';

@Component({
    selector: 'app-image',
    templateUrl: './image.component.html',
    styleUrls: ['./image.component.css']
})
export class ImageComponent implements OnInit {
    images: Image[];

    private subscription: Subscription;

    constructor(private imageService: ImageService, private eventService: EventService) {
        this.images = [];
    }

    fetchImages() {
        this.imageService.getImages().then(images => {
            console.log(images);
            this.images = images;
        });
    }

    ngOnInit() {
        this.fetchImages();

        this.subscription = this.eventService.event.filter((event): boolean => {
            return (event.Type == 'image' && (event.Action == 'pull' || event.Action == 'delete'));
        }).subscribe((event: any) => {
            console.log('pull:', event.id);

            this.fetchImages();
        });
    }

    ngOnDestroy() {
        this.subscription.unsubscribe();
    }

    pull(image: Image) {
        console.log(image);

        this.imageService.pull(image.Repository);
    }

    remove(image: Image) {
        this.imageService.remove(image.Id);
    }

}
