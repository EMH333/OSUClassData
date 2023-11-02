import typescript from "svelte-preprocess";

const options  = {
    preprocess: typescript({}),
    compilerOptions: {
        hydratable: true,
    }
};

export default options;
