<template>
	<div class='container'>
		<Button @click="refresh()" label='Refresh' color="black"/>
		<AddRoom @add-room="addRoom"/>
		<div class='room-container'>
			<div :key='room.name'  v-for='room in rooms'>
			<Room :room="room"/>	
			</div>
		</div>
	</div>
</template>


<script>
import Button from './Button';
import Room from './Room';
import AddRoom from './AddRoom';

export default{
	name: 'Rooms',
	components:{
		Button, Room, AddRoom
	},
	props:{
		rooms: Array,
	},
	methods: {
		refresh(){
			this.$emit('refresh-rooms')
		},
		addRoom(room){
			this.$emit('add-Room', room)
		}
	}

}
</script>

<style scoped>
.container {
  max-width: 300px;
  margin: 20px;
  overflow: auto;
  min-height: 300px;
  border: 1px solid steelblue;
  padding: 20px;
  border-radius: 5px;
}
.room-container {
  max-width: 300px;
  margin: 3px;
  overflow: auto;
  min-height: 300px;
  border: 1px solid steelblue;
  padding: 10px;
  border-radius: 5px;
}

</style>

