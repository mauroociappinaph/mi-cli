#!/bin/bash
set -e

# Sync documentation and SDD artifacts to Notion
# Requires NOTION_TOKEN and NOTION_DATABASE_ID environment variables

if [ -z "$NOTION_TOKEN" ] || [ -z "$NOTION_DATABASE_ID" ]; then
    echo "⚠️  NOTION_TOKEN or NOTION_DATABASE_ID not set, skipping Notion sync"
    exit 0
fi

echo "📤 Syncing to Notion..."

# Use opencode with Notion MCP to sync
cat > /tmp/notion-sync-prompt.md << 'PROMPTEOF'
Sync the following to Notion database:

1. **README.md** - Project overview
2. **CHANGELOG.md** - Recent changes
3. **API.md** - API reference
5. **SDD Artifacts** (.atl/):
   - Proposals: .atl/proposals/
   - Specs: .atl/specs/
   - Designs: .atl/designs/
   - Tasks: .atl/tasks/

For each file, create/update a Notion page in the database with:
- Title: filename
- Content: file content (markdown)
- Tags: type (readme, changelog, api, proposal, spec, design, task)
- Project: Ayrton CLI
- Last synced: timestamp
PROMPTEOF

# Run via opencode with Notion MCP
opencode run --agent orchestrator-agent "$(cat /tmp/notion-sync-prompt.md)" || {
    echo "⚠️  Notion sync via opencode failed, trying direct API..."
    # Fallback: direct Notion API would go here
    echo "ℹ️  Manual Notion sync needed - artifacts in .atl/ and docs/"
}

echo "✅ Notion sync attempted"
