<script setup lang="tsx">
const columns = [
  {
    name: "participantNumber",
    required: true,
    label: "Id",
    align: "left",
    field: (row) => row.participantNumber,
    format: (val) => `${val}`,
    sortable: true,
  },
  {
    name: "name",
    required: true,
    label: "Name",
    align: "left",
    field: (row) => row.name,
    format: (val) => `${val}`,
    sortable: true,
  },
  { name: "actions", 
    label: "Actions", 
    field: "",
    align: "left" 
    },
];

 const pagination = ref({
      sortBy: 'asc',
      descending: false,
      page: 1,
      rowsPerPage: 12
      // rowsNumber: xx if getting data from a server
    })

interface Participant {
    name: string;
}

</script>

<script lang="tsx">
  import { resolveComponent, defineComponent, reactive, PropType, ref, onMounted } from 'vue';
  import { axiosApiInstance } from '../main';
  import { QTable } from 'quasar'
 
export default defineComponent ({
//  export default {
        name: 'Participants',
        inject: ['router'],

        data: function() {
            return {
                loading: true,
                participants: [],
///                tournament: {id:0} as Tournament
//                tournament: {}
            }
        },
        watch: {
            currentTab(val) {
                if (val == 'Participants') {
                    console.log("Participants activated")
                }
            }
        },

        props: {
           currentTab: String,
          tournamentId: Number,
        },

        mounted() {

            this.fetchParticipants()
        },

        updated() {
           // this.fetchParticipants()
        },

        methods: {

          fetchParticipants() {
//            const tableComp = resolveComponent('MyComponent') as QTable

            this.loading = true
            let cmp = this;
            if (cmp.tournamentId == 0) {
                return
            }
            console.log("cmp.tournamentId = " + JSON.stringify(cmp.tournamentId))
            let url = "/tournaments/" + cmp.tournamentId + "/participants"
            axiosApiInstance.get(url)
                .then(function (response) {
                    console.log("response = " + JSON.stringify(response))
                    cmp.participants = response.data
                   // tableComp.sort("participantNumber")
    //p.table.rows = response.data
                })
                .catch(function (error) {
                    // handle error
                     console.log(error);
                })
                .then(function () {
                      cmp.loading = false
                     // always executed
                });
            },
            AddPlayer() {
              var txt = document.getElementById('participant_names') as HTMLTextAreaElement;
              let usernames = txt.value.split(",");
              var namearray = new Array()
              let cmp = this
              for (let nm of usernames) {
                 if (nm) {
                    let n = nm.trim();
                    if (n.length > 0) {
                        let name = {name: nm}
                        namearray.push(name)
                    }
                 }
              }
              if (namearray.length <= 0) {
                return
              }

              let url = "/tournaments/" + cmp.tournamentId + "/participants"
              axiosApiInstance.post(url , JSON.stringify(namearray))
                  .then(function (response) {
                     console.log("response = " + JSON.stringify(response))
                     cmp.participants = response.data
                     cmp.showInfoBanner()

                  })
                  .catch(function (error) {
                     // handle error
                     console.log(error);
                     cmp.toggleErrorBanner()
                  }).then(function () {
                     // always executed
                  });

            },
            deleteParticipant(props) {

                console.log("props: "+ JSON.stringify(props))

                let cmp = this
                let url = "/tournaments/" +cmp.tournamentId + "/participants/" + props.row.participantNumber;
                   
                axiosApiInstance.delete(url)
                    .then(function (response) {

                        console.log("response = " + JSON.stringify(response))
                        cmp.participants = response.data

                    })
                    .catch(function (error) {
                        // handle error
                        console.log(error);
                    }).then(function () {
                    })
            },
            setPaid(props) {
                console.log("props: "+ JSON.stringify(props))

                let cmp = this
                let url = "/tournaments/" + cmp.tournamentId + "/participants/" + props.row.participantNumber + 
                   "/paid"

                axiosApiInstance.post(url)
                    .then(function (response) {
                        console.log("response = " + JSON.stringify(response))
                        cmp.participants = response.data
                    })
                    .catch(function (error) {
                        // handle error
                        console.log(error);
                    }).then(function () {
                        // always executed
                    })

            },
            saveName(props) {
                console.log("value = " + JSON.stringify(props))
                let cmp = this
                const row = props["row"]
                console.log("row = " + JSON.stringify(row))
                const pNumber = row["participantNumber"]

                let url = "/tournaments/" + cmp.tournamentId + "/participants/" + pNumber
                const name = props["key"]

                axiosApiInstance.put(url, name)
                    .then(function (response) {
                        console.log("response = " + JSON.stringify(response))
                    })
                    .catch(function (error) {
                        // handle error
                        console.log(error);
                    }).then(function () {
                        // always executed
                    })
            }, 
            showInfoBanner() {
              console.log("showInfoBanner")
              var infoBanner = document.getElementById("info-banner");
              var infoBanner = document.getElementById("info-banner");
              var errorBanner = document.getElementById("error-banner");
              if (errorBanner != null) {
                  errorBanner.style.display = "none";
              }
              if (infoBanner != null) {
                  infoBanner.innerHTML="Player(s) added"
                  infoBanner.style.display = "block";
                  infoBanner.setAttribute("class", "bg-positive text-white")
                  infoBanner.setAttribute("dense", "true")
              }
              var txt = document.getElementById('participant_names') as HTMLTextAreaElement;
              txt.value = ""

            },
            toggleErrorBanner() {
              console.log("togglErrrorBanner")
              var infoBanner = document.getElementById("info-banner");
              var errorBanner = document.getElementById("error-banner");
              if (infoBanner != null) {
                  infoBanner.style.display = "none";
              }
              if (errorBanner != null) {
                if (errorBanner.style.display === "none") {
                  errorBanner.innerHTML="No more new players can be accepted at this point"
                  errorBanner.setAttribute("class", "bg-negative text-white")
                  errorBanner.setAttribute("dense", "false")
                  errorBanner.style.display = "block";
                } else {
                  errorBanner.style.display = "none";
                }
              }
            }
            /*
            ,
            updateList() {
                const table = document.querySelector("#tournament_list");
                for (const currentRow of table.rows) {
                     currentRow.onclick = this.createClickHandler(currentRow);
                }
             }
             */
        }
    })

</script>

<template>
  <div v-if="loading">
    <h1>Loading ....</h1>
  </div>

  <div id="q-app" style="min-height: 100vh" v-else>
    <q-banner dense id="info-banner" style="display:none" class="bg-positive text-white">
    </q-banner>
    <q-banner dense id="error-banner" style="display:none" class="bg-negative text-white">
      No more new players can be accepted at this point
      <template v-slot:action>
        <q-btn flat color="white" label="Dismiss"
                @click="toggleErrorBanner()"></q-btn>
      </template>
    </q-banner>
    <input type="text" id="participant_names" name="names" />
    <q-btn
      style="display: inline; margin-left: 6px"
      v-on:click="AddPlayer"
      push
      color="primary"
      label="Add"
    ></q-btn>
    <div class="q-pa-md">
      <q-table
        :dense="true"
        title="Players"
        :rows-per-page-option=10
        :rows="participants"
         v-model:pagination="pagination"
        :columns="columns"
        row-key="name"
        id="participant-table"
        ref="tableRef"
      >
        <template v-slot:body="props">
          <q-tr :props="props">
              <q-td key="participantNumber" :props="props">
                {{ props.row.participantNumber }}
              </q-td>
              <q-td key="name" :props="props">
                {{ props.row.name }}
                <q-popup-edit  buttons 
                v-model="props.row.name" auto-save v-slot="scope">
                    <q-input v-model="props.row.name" @keyup.enter="scope.set"
                     dense autofocus @blur="saveName(props)" ></q-input>
                    <!--<q-input v-model="props.row.name" @keyup.enter="scope.set"
                     dense autofocus @update:modelValue="saveName(props)"></q-input>-->
                </q-popup-edit>
              </q-td>
              <q-td key="actions">
                <q-btn
                  v-if="props.row.paidAmount == 0"
                  dense
                  round
                  flat
                  color="red"
                  @click="setPaid(props)"
                  icon="paid"
                ></q-btn>
                <q-btn v-else dense round flat color="green" icon="paid"></q-btn>
                <q-btn
                  dense
                  round
                  flat
                  color="red"
                  @click="deleteParticipant(props)"
                  icon="delete"
                ></q-btn>
              </q-td>
            </q-tr>
            </template>
      </q-table>
    </div>
  </div>
</template>


<style scoped>
.hoverTable {
  width: 100%;
  border-collapse: collapse;
}
.hoverTable td {
  padding: 7px;
  border: #4e95f4 1px solid;
}
/* Define the default color for all the table rows */
.hoverTable tr {
  background: #b8d1f3;
}
/* Define the hover highlight color for the table row */
.hoverTable tr:hover {
  background-color: #ffff99;
}
</style>
