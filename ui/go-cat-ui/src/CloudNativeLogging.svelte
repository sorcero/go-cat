<script lang="ts">
    import type infra_json from "./infra.json";
    export let infra: typeof infra_json["infra"][0];
</script>


{#if infra.infra_type == 'run.googleapis.com'}
<div class="control">
    <a href="https://console.cloud.google.com/run/detail/us-east1/{infra.name}/logs?project={infra.cloud_project_id}" target="_blank">
        <div class="tags has-addons">
            <span class="tag is-dark">Logging</span>
            <span class="tag is-primary">Google</span>
        </div>
    </a>
</div>
{/if}