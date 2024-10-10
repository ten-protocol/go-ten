// @ts-expect-error
import { defineConfig } from "@openapi-codegen/cli";
import {
  generateReactQueryComponents,
  generateSchemaTypes,
} from "@openapi-codegen/typescript";

export default defineConfig({
  digest: {
    from: {
      relativePath: "./storage/digest-openapi.json",
      source: "file",
    },
    outputDir: "./src/api/codegen",
    // @ts-expect-error
    to: async (context) => {
      const filenamePrefix = "digest";

      const { schemasFiles } = await generateSchemaTypes(context, {
        filenamePrefix,
      });
      return await generateReactQueryComponents(context, {
        filenamePrefix,
        schemasFiles,
      });
    },
  },
});
