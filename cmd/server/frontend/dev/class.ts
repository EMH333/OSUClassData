import Class from './class.svelte';
import { mount } from "svelte";

document.addEventListener('DOMContentLoaded', function () {
    mount(Class, {
            target: document.body
        })
});
