<template>
		<button @click="onClick()" :style="{ background: color }" class='btn'>{{ label }}</button>
</template>


<script>
	export default{
		name: 'Button',
		props:{
			label: String,
			color: String
		},
		methods:{
			onClick(){
				console.log('click')
			}
		}
	}
</script>

