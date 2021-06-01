package client

import (
	"log"

	"github.com/calamity-of-subterfuge/cos/v2/pkg/srvpkts"
	"github.com/jakecoffman/cp"
	"github.com/mitchellh/mapstructure"
)

// GameObject describes some object that's within the game. This is often
// wrapped.
type GameObject struct {
	// UID of the object
	UID string

	// SheetURL describes where the PixiJS spritesheet file containing
	// the animations for this object can be found
	SheetURL string

	// SpriteScale describes the x/y scaling for this game object compared
	// to the underlying sprite.
	SpriteScale cp.Vector

	// SpriteRotation is the default rotation in radians for this game object
	// compared to how it is on the spritesheet, which can be used for object
	// packing
	SpriteRotation float64

	// RenderOffset is the visual offset of the sprite on the screen in pixels
	// at standard zoom.
	RenderOffset cp.Vector

	// Animation is the name of the animation within the sheet that
	// should currently be used to render this object.
	Animation string

	// AnimationSpeed is the speed of the animation, where 1 is 60
	// frames per second, 0.5 is 30 frames per second, and 2 is
	// 120 frames per second.
	AnimationSpeed float64

	// AnimationPlaying is true if the animation should actually play
	// and false if it should be frozen on an arbitrary frame.
	AnimationPlaying bool

	// AnimationLooping is true if the animation should loop when
	// finished and false if it should not.
	AnimationLooping bool

	// Body is the physics body for this game object, with the appropriate
	// shapes, position, velocity, angle, and angular velocity, but not
	// attached to any space and hence not simulated.
	Body *cp.Body
}

// Sync this game object using the given information. This overwrites
// everything
func (o *GameObject) Sync(packet *srvpkts.GameObjectSync) *GameObject {
	o.UID = packet.UID
	o.SheetURL = packet.SheetURL
	o.SpriteScale = cp.Vector{X: packet.SpriteScale.X, Y: packet.SpriteScale.Y}
	o.SpriteRotation = packet.SpriteRotation
	o.RenderOffset = cp.Vector{X: packet.RenderOffset.X, Y: packet.RenderOffset.Y}
	o.Animation = packet.Animation
	o.AnimationSpeed = packet.AnimationSpeed
	o.AnimationPlaying = packet.AnimationPlaying
	o.AnimationLooping = packet.AnimationLooping

	body := cp.NewBody(0, 0)
	for _, shp := range packet.Shapes {
		body.AddShape(makeCPShape(body, &shp))
	}

	body.SetPosition(cp.Vector{X: packet.Position.X, Y: packet.Position.Y})
	body.SetVelocity(packet.Velocity.X, packet.Velocity.Y)
	body.SetAngle(packet.Rotation)
	body.SetAngularVelocity(packet.AngularVelocity)
	return o
}

// Update this game object with the given information, which only affects
// the relatively frequently changing fields
func (o *GameObject) Update(packet *srvpkts.GameObjectUpdatePacket) *GameObject {
	o.Body.SetPosition(cp.Vector{X: packet.Position.X, Y: packet.Position.Y})
	o.Body.SetVelocity(packet.Velocity.X, packet.Velocity.Y)
	o.Body.SetAngle(packet.Rotation)
	o.Body.SetAngularVelocity(packet.AngularVelocity)
	o.Animation = packet.Animation
	o.AnimationPlaying = packet.AnimationPlaying
	o.AnimationLooping = packet.AnimationLooping
	return o
}

func makeCPShape(body *cp.Body, shp *srvpkts.Shape) *cp.Shape {
	switch shp.ShapeType {
	case "polygon":
		var details srvpkts.PolygonDetails
		if err := mapstructure.Decode(shp.Details, &details); err != nil {
			log.Fatalf("error decoding polygon details in %v: %v", shp, err)
		}

		cpVerts := make([]cp.Vector, len(details.Vertices))
		for idx, vert := range details.Vertices {
			cpVerts[idx] = cp.Vector{X: vert.X, Y: vert.Y}
		}

		res := cp.NewPolyShape(body, len(details.Vertices), cpVerts, cp.NewTransformIdentity(), details.Radius)
		res.SetMass(shp.Mass)
		return res
	default:
		log.Fatalf("unknown shape type: %v", shp.ShapeType)
		return nil
	}
}
