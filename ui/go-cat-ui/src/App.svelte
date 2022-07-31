<script lang="ts">
  // utils: relative time and copy to clipboard
  import RelativeTime from '@yaireo/relative-time'
  import clipboardCopy from 'clipboard-copy'

  // search
  import Fuse from 'fuse.js'

  // svelte components
  import CloudNativeMonitoring from './CloudNativeMonitoring.svelte';
  import CloudNativeLogging from './CloudNativeLogging.svelte';

  // add fontawesome icons 
  import Fa from 'svelte-fa/src/fa.svelte'
  import { faClipboard, faArrowUpRightFromSquare, faLink } from '@fortawesome/free-solid-svg-icons'
  
  // flag which will be set to true when the user is currently 
  // typing something in the input box
  let searching = false;

  // get relative time object for calulating relative time according
  // to when the environment list was updated, and also for deployment 
  // time
  const relativeTime = new RelativeTime(); 
	
  import infra_json from './infra.json';
import Monitoring from './Monitoring.svelte';
import InfraType from './InfraType.svelte';

  let data = infra_json["infra"];

  const fuse = new Fuse(data, {
    keys: [
      "id", 
      {
        name: "cloud_project_id",
        weight: 2,
      }, 
      {
        name: "name",
        weight: 4
      },
      "subsystem"
    ]
  })

  let search_timer: ReturnType<typeof setTimeout> | null = null;
  let search_input_value = "";
  function triggerUpdate() {
    clearTimeout(search_timer)

    search_timer = setTimeout(function() {

      if (search_input_value == "") {
        data = infra_json["infra"];
        searching = false;
        return
      }
      
      let temp_data = fuse.search(search_input_value)
      data = temp_data.map(({ item })  => item);

      searching = false;
      //console.log(data)
    }, 700)
    
  }

  

</script>



<main>
<section class="section">
  <div class="container">
  
    <h1 class="title">Environment List</h1>
    <p id="last_updated" class="subtitle is-loading">
        Last updated <strong>{relativeTime.from(new Date(infra_json["updated_at"]))}</strong>. 
    </p>

    <div class="control {searching === true ? "is-loading": ""}">
      <input class="input is-primary" type="text" placeholder="Quick search" on:input={() => {
        searching = true
        triggerUpdate()
      }} bind:value="{search_input_value}">
    </div>

    <br>

    <table class="table is-fullwidth">
      <thead>
          <tr>
              <th><abbr title="Deployment Name">Name</abbr></th>
              <th>Deployment Links</th>
              <th>Tags</th>
          </tr>

      </thead>
      <tbody id="gocat__cloud_body">
        
        {#each data as d}
        <tr>
            <td>
              <a id="{d.id}" href="#{d.id}" class="icon-text is-small">
                <span class="icon"><Fa icon={faLink} style="font-size: 0.75em"  /></span>
              </a> 
              {d.cloud.toLowerCase()} > {d.cloud_project_id} >
              <br>
              {d.subsystem}
              > <br>
              <strong class="is-underlined">
                
                {d.name}
              </strong>
              <button class="button is-small is-text" on:click={() => {clipboardCopy(d.id)} }>
                <span class="icon is-small">
                  <Fa icon={faClipboard}/>
                </span>
              </button>
              <br>
              
          </td>	
            {#if d.deployment_links} 
            <td style="word-break: break-all;">								
            {#each d.deployment_links as link}	
              <a href="{link}">{link}</a>
              <button class="button is-small is-text" on:click={() => {clipboardCopy(link)}}>
                <span class="icon is-small">
                  <Fa icon={faClipboard} />
                </span>
              </button>
              <a class="button is-small is-text" target="_blank" href="{link}" >
                <span class="icon is-small">
                  <Fa icon={faArrowUpRightFromSquare} />
                </span>
              </a>
              <br>
              {/each}
            </td>
            {:else if d.deployment_link}
            <td>
              <a href="{d.deployment_link}">{d.deployment_link}</a>
              <button class="button is-small is-text" on:click={() => {clipboardCopy(d.deployment_link)}}>
                <span class="icon is-small">
                  <Fa icon={faClipboard} />
                </span>
              </button>
            </td>
          {:else}
          <td>No deployment link</td>
          {/if}

            

          
            
          <td>
            <div class="field is-grouped is-grouped-multiline">
              <div class="control">
                <div class="tags has-addons">
                  <span class="tag is-dark">Cloud</span>
                  <span class="tag is-info">{d.cloud}</span>
                </div>
              </div>

              <div class="control">
                <div class="tags has-addons">
                  <span class="tag is-dark">Project</span>
                  <span class="tag is-link">{d.cloud_project_id}</span>
                </div>
              </div>
              
              <div class="control">
                <div class="tags has-addons">
                  <span class="tag is-dark">Deployed</span>
                  <span class="tag is-light"><time datetime={d.deployed_on}>{relativeTime.from(new Date(d.deployed_on))}</time></span>
                </div>
              </div>

              <div class="control">
                <div class="tags has-addons">
                  <span class="tag is-dark">SHA</span>
                  <span class="tag is-primary">{d.commit_sha}</span>
                </div>
              </div>

              {#if d.infra_type }
              <InfraType infra={d} />
              <CloudNativeMonitoring infra={d} />
              <CloudNativeLogging infra={d} />
              {/if}

              {#if d.monitoring_links}
                <Monitoring infra={d} />
              {/if}
            </div>
          
          </td>
        </tr>
        {/each}
        
      </tbody>
    </table>

  </div>
</section>	
</main>