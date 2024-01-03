package engine

import (
	"unsafe"
	gl "github.com/chsc/gogl/gl33"
)

const (
	// main shaders
	vertexShaderSource = `
	    #version 330
		attribute vec3 pos;
		attribute vec2 uv;
		attribute vec4 color;

		varying vec4 v_color;
		varying vec2 v_uv;
		uniform mat4 view;
		uniform mat4 model;
		uniform mat4 projection;
		uniform vec2 screen;
		uniform vec3 camera_pos;
		uniform vec2 fade;
		uniform float time;
		
		void main() {
			gl_Position = projection * view * model * vec4(pos, 1.0);
			gl_Position.xy += screen.xy * gl_Position.w;
			v_color = color;
			v_color.a *= smoothstep(
				fade.y, fade.x, // fadeout far, near
				length(vec4(camera_pos, 1.0) - model * vec4(pos, 1.0))
			);
			v_uv = uv / 2048.0; // ATLAS_GRID * ATLAS_SIZE
		}
	`

	fragmentShaderSource = `
	    #version 330
		varying vec4 v_color;
		varying vec2 v_uv;
		uniform sampler2D texture;

		void main() {
			vec4 tex_color = texture2D(texture, v_uv);
			vec4 color = tex_color * v_color;
			if (color.a == 0.0) {
				discard;
			}
			color.rgb = color.rgb * 2.0;
			gl_FragColor = color;
		}
	`
)

func createProgram(vsSource string, fsSource string) gl.Uint {
	// VERTEX SHADER
	vs := compileShader(gl.VERTEX_SHADER, vsSource)
	fs := compileShader(gl.FRAGMENT_SHADER, fsSource)

	// CREATE PROGRAM
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.UseProgram(program)
	// fragoutstring := gl.GLString("outColor")
	// defer gl.GLStringFree(fragoutstring)
	// gl.BindFragDataLocation(program, gl.Uint(0), fragoutstring)

	// var linkstatus gl.Int
	// gl.GetProgramiv(program, gl.LINK_STATUS, &linkstatus)
	// fmt.Printf("Program Link: %v\n", linkstatus)

	return program
}

func compileShader(shaderType int, source string) gl.Uint{
	shader := gl.CreateShader(gl.Enum(shaderType))
    ssource := gl.GLString(source)
	gl.ShaderSource(shader, 1, &ssource, nil)
	gl.CompileShader(shader)

	var success gl.Int
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success != 1 {
	var logLength gl.Int 
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		var log string
		gl.GetShaderInfoLog(shader, gl.Sizei(logLength), nil, gl.GLString(log))
		Logger.Printf("Compiled Shader %d: %v\n", shaderType, success)
	}

	return shader
}

type uniform struct {
	view       gl.Uint
	model      gl.Uint
	projection gl.Uint
	screen     gl.Uint
	cameraPos  gl.Uint
	fade       gl.Uint
	time       gl.Uint
}

type attribute struct {
	pos   gl.Uint
	uv    gl.Uint
	color gl.Uint
}

type ProgramGame struct {
	program   gl.Uint
	vao       gl.Uint
	uniform   uniform
	attribute attribute
}

// func bindVaF(index gl.Uint, container, member, start) {
// 	gl.VertexAttribPointer(index, 3, gl.FLOAT, gl.FALSE, 9*4, gl.Pointer(nil))
// }

// func bindVaColor(index gl.Uint, container, member, start) {
// 	gl.VertexAttribPointer(index, 3, gl.FLOAT, gl.FALSE, 9*4, gl.Pointer(nil))
// }

func ShaderGameInit() *ProgramGame {
	p := &ProgramGame{}
	p.program = createProgram(vertexShaderSource, fragmentShaderSource)

	p.uniform.view = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("view")))
	p.uniform.model = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("model")))
	p.uniform.projection = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("projection")))
	p.uniform.screen = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("screen")))
	p.uniform.cameraPos = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("camera_pos")))
	p.uniform.fade = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("fade")))
	p.uniform.time = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("time")))

	p.attribute.pos = gl.Uint(gl.GetAttribLocation(p.program, gl.GLString("pos")))
	p.attribute.uv = gl.Uint(gl.GetAttribLocation(p.program, gl.GLString("uv")))
	p.attribute.color = gl.Uint(gl.GetAttribLocation(p.program, gl.GLString("color")))

	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	gl.EnableVertexAttribArray(p.attribute.pos)
	gl.EnableVertexAttribArray(p.attribute.uv)
	gl.EnableVertexAttribArray(p.attribute.color)

	// glVertexAttribPointer( s->attribute.pos, 
	//  sizeof(((vertex_t *)0)->pos)/sizeof(float), 
	//  0x1406, 
	//  0, 
	//  sizeof(vertex_t), 
	//  (GLvoid*)(__builtin_offsetof (vertex_t, pos) + 0) )
	gl.VertexAttribPointer(p.attribute.pos, 3, gl.FLOAT, gl.FALSE, gl.Sizei(unsafe.Sizeof(Vertex{})), gl.Pointer(uintptr(unsafe.Offsetof(Vertex{}.Pos))))
	gl.VertexAttribPointer(p.attribute.uv, 2, gl.FLOAT, gl.FALSE, gl.Sizei(unsafe.Sizeof(Vertex{})), gl.Pointer(uintptr(unsafe.Offsetof(Vertex{}.UV))))
	gl.VertexAttribPointer(p.attribute.color, 4, gl.UNSIGNED_BYTE, gl.FALSE, gl.Sizei(unsafe.Sizeof(Vertex{})), gl.Pointer(uintptr(unsafe.Offsetof(Vertex{}.Color))))

	return p
}

// Post effect shaders
const (
	postEffectVertexShaderSource = `
		attribute vec3 pos;
	attribute vec2 uv;

	varying vec2 v_uv;

	uniform mat4 projection;
	uniform vec2 screen_size;
	uniform float time;
	
	void main() {
		gl_Position = projection * vec4(pos, 1.0);
		v_uv = uv;
	}`

	postEffectFragmentShaderSourceDefault = `
	varying vec2 v_uv;

	uniform sampler2D texture;
	uniform vec2 screen_size;

	void main() {
		gl_FragColor = texture2D(texture, v_uv);
	}
	`
	postEffectFragmentShaderSourceCRT = `
		varying vec2 v_uv;

	uniform float time;
	uniform sampler2D texture;
	uniform vec2 screen_size;

	vec2 curve(vec2 uv) {
		uv = (uv - 0.5) * 2.0;
		uv *= 1.1;	
		uv.x *= 1.0 + pow((abs(uv.y) / 5.0), 2.0);
		uv.y *= 1.0 + pow((abs(uv.x) / 4.0), 2.0);
		uv  = (uv / 2.0) + 0.5;
		uv =  uv *0.92 + 0.04;
		return uv;
	}

	void main(){
		vec2 uv = curve(gl_FragCoord.xy / screen_size);
		vec3 color;
		float x =  sin(0.3*time+uv.y*21.0)*sin(0.7*time+uv.y*29.0)*sin(0.3+0.33*time+uv.y*31.0)*0.0017;

		color.r = texture2D(texture, vec2(x+uv.x+0.001,uv.y+0.001)).x+0.05;
		color.g = texture2D(texture, vec2(x+uv.x+0.000,uv.y-0.002)).y+0.05;
		color.b = texture2D(texture, vec2(x+uv.x-0.002,uv.y+0.000)).z+0.05;
		color.r += 0.08*texture2D(texture, 0.75*vec2(x+0.025, -0.027)+vec2(uv.x+0.001,uv.y+0.001)).x;
		color.g += 0.05*texture2D(texture, 0.75*vec2(x+-0.022, -0.02)+vec2(uv.x+0.000,uv.y-0.002)).y;
		color.b += 0.08*texture2D(texture, 0.75*vec2(x+-0.02, -0.018)+vec2(uv.x-0.002,uv.y+0.000)).z;

		color = clamp(color*0.6+0.4*color*color*1.0,0.0,1.0);

		float vignette = (0.0 + 1.0*16.0*uv.x*uv.y*(1.0-uv.x)*(1.0-uv.y));
		color *= vec3(pow(vignette, 0.25));

		color *= vec3(0.95,1.05,0.95);
		color *= 2.8;

		float scanlines = clamp( 0.35+0.35*sin(3.5*time+uv.y*screen_size.y*1.5), 0.0, 1.0);
		
		float s = pow(scanlines,1.7);
		color = color * vec3(0.4+0.7*s);

		color *= 1.0+0.01*sin(110.0*time);
		if (uv.x < 0.0 || uv.x > 1.0)
			color *= 0.0;
		if (uv.y < 0.0 || uv.y > 1.0)
			color *= 0.0;
		
		color*=1.0-0.65*vec3(clamp((mod(gl_FragCoord.x, 2.0)-1.0)*2.0,0.0,1.0));
		gl_FragColor = vec4(color,1.0);
	}
	`
)

type uniformPE struct {
	projection gl.Uint
	screenSize gl.Uint
	time       gl.Uint
}

type attributePE struct {
	pos gl.Uint
	uv  gl.Uint
}

type ProgramPostEffect struct {
	program   gl.Uint
	vao       gl.Uint
	uniform   uniformPE
	attribute attributePE
}

func ShaderPostEffectGeneralInit(p *ProgramPostEffect) {
	p.uniform.projection = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("projection")))
	p.uniform.screenSize = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("screen_size")))
	p.uniform.time = gl.Uint(gl.GetUniformLocation(p.program, gl.GLString("time")))

	p.attribute.pos = gl.Uint(gl.GetAttribLocation(p.program, gl.GLString("pos")))
	p.attribute.uv = gl.Uint(gl.GetAttribLocation(p.program, gl.GLString("uv")))

	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	gl.EnableVertexAttribArray(p.attribute.pos)
	gl.EnableVertexAttribArray(p.attribute.uv)

	gl.VertexAttribPointer(p.attribute.pos, 3, gl.FLOAT, gl.FALSE, gl.Sizei(unsafe.Sizeof(Vertex{})), gl.Pointer(uintptr(unsafe.Offsetof(Vertex{}.Pos))))
	gl.VertexAttribPointer(p.attribute.uv, 2, gl.FLOAT, gl.FALSE, gl.Sizei(unsafe.Sizeof(Vertex{})), gl.Pointer(uintptr(unsafe.Offsetof(Vertex{}.UV))))
}

func ShaderPostEffectDefaultInit() *ProgramPostEffect {
	p := &ProgramPostEffect{}
	p.program = createProgram(postEffectVertexShaderSource, postEffectFragmentShaderSourceDefault)
	ShaderPostEffectGeneralInit(p)

	return p
}

func ShaderPostEffectCRTInit() *ProgramPostEffect {
	p := &ProgramPostEffect{}
	p.program = createProgram(postEffectVertexShaderSource, postEffectFragmentShaderSourceCRT)
	ShaderPostEffectGeneralInit(p)

	return p
}
