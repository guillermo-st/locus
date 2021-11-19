<template>
	<h1> Locus. </h1>
	<Rooms :rooms="rooms" @refresh-rooms="getRooms()"/>
</template>

<script>
import Rooms from './components/Rooms';
import Button from './components/Button';


export default {
	name: 'App',
  	components: {
  		Button,Rooms
  	},
  	data(){
		return{
			rooms: null
		}
  	},
	methods: {
		getRooms(){
			fetch('http://localhost:8000/rooms')
				.then(response => response.json())
				.then(data => (this.rooms = data))
		}
	},
 	created() {
		this.getRooms()
  	}
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400&display=swap');
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
body {
  font-family: 'Poppins', sans-serif;
}
.container {
  max-width: 300px;
  margin: 20px auto;
  overflow: auto;
  min-height: 600px;
  border: 1px solid steelblue;
  padding: 30px;
  border-radius: 5px;
}
.btn {
  display: inline-block;
  background: #000;
  color: #fff;
  border: none;
  padding: 10px 20px;
  margin: 5px;
  border-radius: 5px;
  cursor: pointer;
  text-decoration: none;
  font-size: 15px;
  font-family: inherit;
}
.btn:focus {
  outline: none;
}
.btn:active {
  transform: scale(0.98);
}
.btn-block {
  display: block;
  width: 100%;
}
</style>
