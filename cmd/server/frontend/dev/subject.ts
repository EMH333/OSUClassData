import Subject from './subject.svelte';
import { mount } from "svelte";

document.addEventListener('DOMContentLoaded', function () {
    mount(Subject, {
            target: document.body
        })
});
