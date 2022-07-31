<script lang="ts">
    import type infra_json from "./infra.json";
    export let infra: typeof infra_json["infra"][0];
</script>



{#each infra.monitoring_links as link}
{#if link}
<div class="control">
<a href="{link}" target="_blank">
    <div class="tags has-addons">
    <span class="tag is-dark">Monitoring</span>
    <span class="tag is-primary">{(new URL(link)).hostname}</span>
    </div>
</a>
</div>
{/if}
{/each}