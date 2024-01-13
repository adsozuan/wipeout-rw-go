package main

import (
	"fmt"
	"log"

	"github.com/adsozuan/wipeout-rw-go/engine"
	"github.com/adsozuan/wipeout-rw-go/game"
	"github.com/chsc/gogl/gl33"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Shader sources
var vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;

out vec2 TexCoord;

void main() {
    gl_Position = vec4(aPos, 1.0);
    TexCoord = aTexCoord;
}
` + "\x00"

var fragmentShaderSource = `
#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

uniform sampler2D ourTexture;

void main() {
    FragColor = texture(ourTexture, TexCoord);
}
` + "\x00"

func main() {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("Failed to initialize SDL: %s\n", err)
	}
	defer sdl.Quit()

	// Set OpenGL version (3.3 in this case)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	// Create window
	window, err := sdl.CreateWindow("OpenGL Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	defer window.Destroy()

	// Create OpenGL context
	glContext, err := window.GLCreateContext()
	if err != nil {
		log.Fatalf("Failed to create OpenGL context: %s\n", err)
	}
	defer sdl.GLDeleteContext(glContext)

	// Initialize gogl
	if err := gl33.Init(); err != nil {
		log.Fatalf("Failed to initialize gogl: %s\n", err)
	}

	// Load and create texture
	// texture := loadTexture("C:/Users/ad/wd/wipeout-rw-go/cmd/texture/texture.png")
	texture := loadTimTexture("C:/Users/ad/wd/wipeout-rw-go/cmd/wipeout/data/textures/wiptitle.tim")

	// Compile shaders and create shader program
	shaderProgram := makeShaderProgram(vertexShaderSource, fragmentShaderSource)

	// Set up vertex data and buffers and configure vertex attributes
	var VAO, VBO gl33.Uint
	vertices := []float32{
		// positions   // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, // top right
		0.5, -0.5, 0.0, 1.0, 1.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 1.0, // bottom left
		-0.5, 0.5, 0.0, 0.0, 0.0, // top left
	}
	indices := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	gl33.GenVertexArrays(1, &VAO)
	gl33.GenBuffers(1, &VBO)
	var EBO gl33.Uint
	gl33.GenBuffers(1, &EBO)

	gl33.BindVertexArray(VAO)

	gl33.BindBuffer(gl33.ARRAY_BUFFER, VBO)
	gl33.BufferData(gl33.ARRAY_BUFFER, gl33.Sizeiptr(len(vertices)*4), gl33.Pointer(&vertices[0]), gl33.STATIC_DRAW)

	gl33.BindBuffer(gl33.ELEMENT_ARRAY_BUFFER, EBO)
	gl33.BufferData(gl33.ELEMENT_ARRAY_BUFFER, gl33.Sizeiptr(len(indices)*4), gl33.Pointer(&indices[0]), gl33.STATIC_DRAW)

	// Position attribute
	gl33.VertexAttribPointer(0, 3, gl33.FLOAT, gl33.FALSE, 5*4, nil)
	gl33.EnableVertexAttribArray(0)
	// Texture coord attribute
	gl33.VertexAttribPointer(1, 2, gl33.FLOAT, gl33.FALSE, 5*4, gl33.Pointer(uintptr(3*4)))
	gl33.EnableVertexAttribArray(1)

	// Main loop
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Render
		gl33.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl33.Clear(gl33.COLOR_BUFFER_BIT)

		// Bind texture
		gl33.BindTexture(gl33.TEXTURE_2D, texture)
		var width, height gl33.Int
		gl33.GetTexLevelParameteriv(gl33.TEXTURE_2D, 0, gl33.TEXTURE_WIDTH, &width)
		gl33.GetTexLevelParameteriv(gl33.TEXTURE_2D, 0, gl33.TEXTURE_HEIGHT, &height)

		// Render container
		gl33.UseProgram(shaderProgram)
		gl33.BindVertexArray(VAO)
		gl33.DrawElements(gl33.TRIANGLES, 6, gl33.UNSIGNED_INT, nil)

		// Swap window
		window.GLSwap()

	}

	// Optional: de-allocate all resources once they've outlived their purpose
	gl33.DeleteVertexArrays(1, &VAO)
	gl33.DeleteBuffers(1, &VBO)
	gl33.DeleteProgram(shaderProgram)
}

func loadTimTexture(name string) gl33.Uint {

	data, err := engine.LoadBinaryFile(name)
	if err != nil {
		fmt.Printf("ImageGetTexture-LoadBinaryFile: %s", err)
		return 0
	}
	image := game.ImageLoadFromBytes(data, false)

	// Generate texture
	var texture gl33.Uint
	gl33.GenTextures(1, &texture)
	gl33.BindTexture(gl33.TEXTURE_2D, texture)

	// Set texture parameters
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_WRAP_S, gl33.CLAMP_TO_EDGE)
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_WRAP_T, gl33.CLAMP_TO_EDGE)
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_MIN_FILTER, gl33.LINEAR)
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_MAG_FILTER, gl33.LINEAR)

	// Upload texture to GPU
	gl33.TexImage2D(gl33.TEXTURE_2D, 0, gl33.RGBA, gl33.Sizei(image.Width), gl33.Sizei(image.Height), 0, gl33.RGBA, gl33.UNSIGNED_BYTE, gl33.Pointer(&image.Pixels[0]))

	// Unbind texture
	gl33.BindTexture(gl33.TEXTURE_2D, 0)

	return texture

}

func loadTexture(file string) gl33.Uint {
	img.Init(img.INIT_PNG | img.INIT_JPG)
	defer img.Quit()

	// Load image
	surface, err := img.Load(file)
	if err != nil {
		log.Fatalf("Failed to load image: %s\n", err)
	}
	defer surface.Free()

	// Generate texture
	var texture gl33.Uint
	gl33.GenTextures(1, &texture)
	gl33.BindTexture(gl33.TEXTURE_2D, texture)

	// Set texture parameters
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_WRAP_S, gl33.CLAMP_TO_EDGE)
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_WRAP_T, gl33.CLAMP_TO_EDGE)
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_MIN_FILTER, gl33.LINEAR)
	gl33.TexParameteri(gl33.TEXTURE_2D, gl33.TEXTURE_MAG_FILTER, gl33.LINEAR)

	// Upload texture to GPU
	gl33.TexImage2D(gl33.TEXTURE_2D, 0, gl33.RGBA, gl33.Sizei(surface.W), gl33.Sizei(surface.H), 0, gl33.RGBA, gl33.UNSIGNED_BYTE, gl33.Pointer(&surface.Pixels()[0]))

	// Unbind texture
	gl33.BindTexture(gl33.TEXTURE_2D, 0)

	return texture
}

func makeShaderProgram(vertexSource, fragmentSource string) gl33.Uint {
	vertexShader := gl33.CreateShader(gl33.VERTEX_SHADER)
	cstrVert := gl33.GLString(vertexSource)
	gl33.ShaderSource(vertexShader, 1, &cstrVert, nil)
	gl33.CompileShader(vertexShader)

	fragmentShader := gl33.CreateShader(gl33.FRAGMENT_SHADER)
	cstrFrag := gl33.GLString(fragmentSource)
	gl33.ShaderSource(fragmentShader, 1, &cstrFrag, nil)
	gl33.CompileShader(fragmentShader)

	shaderProgram := gl33.CreateProgram()
	gl33.AttachShader(shaderProgram, vertexShader)
	gl33.AttachShader(shaderProgram, fragmentShader)
	gl33.LinkProgram(shaderProgram)

	gl33.DeleteShader(vertexShader)
	gl33.DeleteShader(fragmentShader)

	return shaderProgram
}
