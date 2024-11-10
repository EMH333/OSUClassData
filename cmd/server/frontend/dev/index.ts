import Index from './index.svelte';
import { hydrate } from "svelte";

document.addEventListener('DOMContentLoaded', function () {
    hydrate(Index, {
            target: document.body
        })
});
