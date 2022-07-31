<script lang="ts">
    import type infra_json from "./infra.json";
    export let infra: typeof infra_json["infra"][0];
</script>

<div class="control">
    <div class="tags has-addons">
        <span class="tag is-dark">Type</span>
        <span class="tag is-success">{infra.infra_type}</span>
    </div>
</div>