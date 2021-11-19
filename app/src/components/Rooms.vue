<template>
	<div class='container'>
		<Button label='Add Room' color="steelblue"/>
		<Button @click="refresh()" label='Refresh' color="black"/>
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

export default{
	name: 'Rooms',
	components:{
		Button, Room
	},
	props:{
		rooms: Array,
	},
	methods: {
		refresh(){
			this.$emit('refresh-rooms')
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

