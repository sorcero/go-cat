<script lang="ts">
    import type infra_json from "./infra.json";
    export let infra: typeof infra_json["infra"][0];

    function togomakURL(infra: typeof infra_json["infra"][0]) {
        const deployedOn: Date = new Date(Date.parse(infra.deployed_on));

        let stage: string = infra.name;
        const overrideStage: any = infra.parameters["togomak.srev.in/v1/stage.id"];
        if (overrideStage && overrideStage !== "") {
            stage = overrideStage;
        }

        let instanceID: string = "";
        const overrideInstanceId: any = infra.parameters["togomak.srev.in/v1/instance.id"];
        if (overrideInstanceId && overrideInstanceId !== "") {
            instanceID = overrideInstanceId;
        }

        let gcl: any = infra.parameters["togomak.srev.in/v1/logging"];
        if (!gcl || (gcl !== "googlde-cloud")) {
            return "";
        }

        let query: string = `jsonPayload.labels.stage = "${stage}"
        labels.instanceId = "${instanceID}"`;
        query = encodeURIComponent(query) + `;timeRange=${deployedOn.toISOString()}/${deployedOn.toISOString()}--PT24H;`;

        return `https://console.cloud.google.com/logs/query;query=${query}?project=${infra.cloud_project_id}`;

    }
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

{#if infra.infra_type == 'togomak.srev.in/release'}
<div class="control">
    <a href="{togomakURL(infra)}" target="_blank">
        <div class="tags has-addons">

            <span class="tag is-dark">
                <i class="fas fa-clipboard"></i>
                Logs
            </span>
            <span class="tag is-primary">Google Cloud Logging</span>
        </div>
    </a>
</div>
{/if}