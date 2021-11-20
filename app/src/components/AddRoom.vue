<template>
  <form @submit="onSubmit" class="add-form">
    <div class="form-control">
      <input type="text" v-model="roomName" name="name" placeholder="Name" maxlength="16" />
    </div>
    <input type="submit" value="Add" class="btn btn-block" />
  </form>
</template>

<script>
export default {
  name: 'AddRoom',
  data() {
    return {
      roomName: '',
    }
  },
  methods: {
    onSubmit(e) {
      e.preventDefault()
      if (!this.roomName) {
        alert("Room name can't be empty!")
        return
      }
      const newRoom = {
        name: this.roomName,
      }
      this.$emit('add-room', newRoom)
      this.name = ''
    },
  },
}
</script>
