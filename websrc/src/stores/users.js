import { defineStore } from 'pinia'
import { reactive } from 'vue'
import { getUsers } from "../api/user";

export const useUsersStore = defineStore('users', () => {
    const users = reactive([])
    const getUsersList = async () => {
        const res = await getUsers()
        users.value = res
    }
    const addUserLocal = (userObj) => {
        users.value.unshift(userObj);
    }
    const delUserLocal = (username) => {
        let indexToRemove = users.value.findIndex(item => item.username === username);
        // console.log(indexToRemove)
        if (indexToRemove !== -1) { // If the index is found (not -1), remove that object from the array
            users.value.splice(indexToRemove, 1);
            // console.log(users.value)
        }
    }
    return { users, getUsersList, addUserLocal, delUserLocal }
})

