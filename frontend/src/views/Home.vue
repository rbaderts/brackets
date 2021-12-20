<script setup lang="ts">

  const columns: any[] = [

    {
      name: "id",
      required: true,
      label: "id",
      align: "left",
      field: (row) => row.id,
      format: (val) => `${val}`,
      sortable: true,
    },
    {
      name: "name",
      align: "center",
      label: "name",
      field: "name",
      sortable: true,
    },
    {
      name: "state",
      align: "center",
      label: "state",
      field: "tournamentState",
      sortable: true,
    },
    {
      name: "playerCcount",
      align: "center",
      label: "state",
      field: "playerCount",
      sortable: true,
    },
  ];

  const pagination = ref({
        sortBy: 'asc',
        descending: false,
        page: 1,
        rowsPerPage: 12
        // rowsNumber: xx if getting data from a server
    })
</script>


<script lang="ts">

  import { defineComponent,ref} from 'vue' 
  import {axiosApiInstance} from '../main.ts';

  export default defineComponent({
  name: "Home",
  inject: ['keycloak', 'router'],

   data() {
    return {
      loading: true,
      tournaments: [{ id: 1, name: "blah" } as Tournament],
      error: null,
      separator: ref("vertical"),
    };
   },  
   mounted: function () {
    let cmp = this;
    this.fetchTournaments().then(function () {
      console.log("tournaments: " + JSON.stringify(cmp.tournaments));
      console.log("ListMounted mounted");
    });
  },

  updated: function () {
    //           this.updateList()
  },
  methods: {
    clickHandler(evt, row) {
      let instance = this;
      console.log("row = " + JSON.stringify(row))
      const id = row['id']
      instance.$router.push({ path: "/tournament/" + id }); // -> /user/123
///      const [cell] = row.getElementsByTagName("td");
 //     const id = cell.innerHTML;
  //    console.log("clickHandler" + id);
//      instance.$router.push({ path: "/tournament/" + id }); // -> /user/123
    },

    createClickHandler(row) {
      let instance = this;
      return () => {
        const [cell] = row.getElementsByTagName("td");
        const id = cell.innerHTML;
        console.log("clickHandler" + id);
        instance.$router.push({ path: "/tournament/" + id }); // -> /user/123
      };
    },

/*
    updateList() {
      console.log("tournaments:  " + JSON.stringify(this.tournaments));
      let table: HTMLTableElement
      table = document.querySelector("#tournament_list table") as HTMLTableElement;

      if (table != null) {
        for (const currentRow of table.rows) {
            currentRow.onclick = this.createClickHandler(currentRow);
        }
      }
    },
    */

    fetchTournaments() {
      //        console.log("keycloak = " + JSON.stringify($this.$keycloak));
      this.error = null;
      this.loading = true;
      let cmp = this;
      return axiosApiInstance
        .get("/tournaments")
        .then(function (response) {
          console.log("response = " + JSON.stringify(response));
          cmp.tournaments = response.data;
          //cmp.updateList();
          //p.table.rows = response.data
        })
        .catch(function (error) {
          // handle error
          cmp.error = error.toString();
          console.log(error);
        })
        .then(function () {
          cmp.loading = false;
          // always executed
        });
    },

    newTournament() {
      console.log("newTourament");
      //        console.log("keycloak = " + JSON.stringify($this.$keycloak));
      let cmp = this;
      return axiosApiInstance
        .post("/tournaments")
        .then(function (response) {
          console.log("response = " + JSON.stringify(response));
          cmp.fetchTournaments().then(function () {
//            cmp.updateList()
          })
          .catch(function(error) {
          })
          //               cmp.table.rows = response.data
        })
        .catch(function (error) {
          // handle error
          cmp.error = error.toString();
          console.log(error);
        })
        .then(function () {
          ///                cmp.isLoading = false
          // always executed
        });
    },
    GetStateString(state) {
      if (state == 1) {
        return "NeedsDraw";
      } else if (state == 2) {
        return "Ready";
      } else if (state == 3) {
        return "Underway";
      } else if (state == 4) {
        return "Registration";
      } else if (state == 5) {
        return "Complete";
      }
    },
  },
});
</script>

<template>
    <div v-if="$keycloak == null || !$keycloak.ready || !$keycloak.authenticated">
        <h1>Not authenticated</h1>
        <button @click="$keycloak.login">Login</button>
    </div>

    <div v-else>

      <q-page>

          <div class="row" style="margin-bottom: 8px">
            <div class="col-2">
              <q-btn
                v-on:click="$keycloak.logoutFn"
                push
                color="primary"
                label="Logout"
              ></q-btn>
            </div>
            <div class="col-2">
              <q-btn
                v-on:click="newTournament"
                push
                color="primary"
                label="New Touranment"
              ></q-btn>
            </div>
          </div>
            <div class="q-pa-md">
              <q-table
                id="tournament_list"
                :dense="true"
                title="Tournaments"
                :rows="tournaments"
                :columns="columns"
                v-model:pagination="pagination"
                @row-click="clickHandler"
                row-key="id" >
                <thead>
                  <tr>
                    <th>Id</th>
                    <th>Name</th>
                    <th>Status</th>
                    <th>Players</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="t in tournaments" :key="t.id">
                    <td>{{ t.id }}</td>
                    <td>{{ t.name }}</td>
                    <td>{{ t.tournamentState }}</td>
                    <td>{{ t.participantCount }}</td>
                  </tr>
                </tbody>
              </q-table>
          <!--
          <h2>You should only be able to see this if you are authenticated.</h2>
          <h3>This is what my token looks like:</h3>

          <button @click="$keycloak.logoutFn">Logout</button>
          <button v-on:click="newTournament">New Tournament</button>
          -->
          

        </div>
      </q-page>
    </div>
  </template>


  <style lang="scss" scoped></style>
