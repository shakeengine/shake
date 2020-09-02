package main

import (
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := sdl.VulkanLoadLibrary(""); err != nil {
		panic(err)
	}
	defer sdl.VulkanUnloadLibrary()

	procAddr := sdl.VulkanGetVkGetInstanceProcAddr()
	if procAddr == nil {
		panic("GetInstanceProcAddr is nil")
	}
	vk.SetGetInstanceProcAddr(procAddr)
	if err := vk.Init(); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("VulkanCube (SDL2)", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(800), int32(600), sdl.WINDOW_VULKAN)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
