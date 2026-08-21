// gen-import.config.js
// Run: npx gen-import --packages --app-config
module.exports = {
     srcDir: [
          'packages',
          "tools"
     ],
     skipPatterns: [
          'packages/index.ts',          // app factory — not part of the barrel
     ],
}
