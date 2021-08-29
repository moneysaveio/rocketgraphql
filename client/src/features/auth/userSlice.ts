import { createSlice } from '@reduxjs/toolkit'

export const userSlice = createSlice({
    name: 'user',
    initialState: {
      value: null
    },
    reducers: {
      set: (state, action) => {
        console.log("Request got: ", action.payload);
        state.value = action.payload
      }
    }
})
  
// Action creators are generated for each case reducer function
console.log("userSlice.actions: ", userSlice.actions);
export const { set } = userSlice.actions

// Other code such as selectors can use the imported `RootState` type
export const selectCount = (state: any) => state.user.value

  
export default userSlice.reducer