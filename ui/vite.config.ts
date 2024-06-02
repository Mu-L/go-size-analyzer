import {defineConfig} from 'vite';
import {viteSingleFile} from "vite-plugin-singlefile"
import * as fs from "node:fs"
import {build, codecov, commonPlugin, getVersionTag} from "./common";
import {createHtmlPlugin} from "vite-plugin-html";


const placeHolder = `"GSA_PACKAGE_DATA"`

const getPlaceHolder = (): string => {
    if (process.env.NODE_ENV === "production") {
        return placeHolder
    }

    try {
        return fs.readFileSync(
            new URL("./src/testdata/testdata.json", import.meta.url),
            "utf-8"
        )
    } catch (e) {
        console.error("Failed to load data.json, for dev you should create one with gsa", e)
        return placeHolder
    }
}

export default defineConfig({
    plugins: [
        ...commonPlugin(),
        createHtmlPlugin({
            minify: true,
            entry: './src/main.tsx',
            inject: {
                tags: [
                    {
                        injectTo: "head",
                        tag: "script",
                        attrs: {
                            type: "application/json",
                            id: "data"
                        },
                        children: getPlaceHolder()
                    },
                    getVersionTag(),
                ]
            }
        }),
        viteSingleFile(
            {
                removeViteModuleLoader: true
            }
        ),
        codecov("gsa-ui"),
    ],
    clearScreen: false,
    build: build("webui"),
})
