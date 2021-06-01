package srvpkts

// GameObjectSync describes everything required for a client to learn about a
// game object its never heard of before.
type GameObjectSync struct {
	// UID of the object
	UID string `mapstructure:"uid" json:"uid"`

	// Sprite

	// SheetURL is the URL to the json object that contains the vast majority of
	// information about how to display this object; it is formatted as a
	// PixiJS-compatible spritesheet with an additional animationMeta section
	// which may include additional render hints not controlled by the server.
	SheetURL string `mapstructure:"sheet_url" json:"sheet_url"`

	// SpriteScale determines the scale that the object should be
	// rendered at. This is in addition to any scaling described by the
	// sheet. This is used to allow dynamic sizing on some objects, e.g.,
	// walls.
	SpriteScale Vector `mapstructure:"sprite_scale" json:"sprite_scale"`

	// SpriteRotation determines a static rotation offset in radians,
	// used to correct for certain types of rotations.
	SpriteRotation float64 `mapstructure:"sprite_rotation" json:"sprite_rotation"`

	// RenderOffset determines a static render offset in pixels, used
	// to correct the rendered bounds to match the physics offset in
	// some situations. Note that the spritesheet may
	RenderOffset Vector `mapstructure:"render_offset" json:"render_offset"`

	// SpriteStateInfo

	// Animation is the name of the animation within the sheet that
	// should currently be used to render this object.
	Animation string `mapstructure:"animation" json:"animation"`

	// AnimationSpeed is the speed of the animation, where 1 is 60
	// frames per second, 0.5 is 30 frames per second, and 2 is
	// 120 frames per second.
	AnimationSpeed float64 `mapstructure:"animation_speed" json:"animation_speed"`

	// AnimationPlaying is true if the animation should actually play
	// and false if it should be frozen on an arbitrary frame.
	AnimationPlaying bool `mapstructure:"animation_playing" json:"animation_playing"`

	// AnimationLooping is true if the animation should loop when
	// finished and false if it should not.
	AnimationLooping bool `mapstructure:"animation_looping" json:"animation_looping"`

	// Body

	// Shapes describes the physics shapes of this body.
	Shapes []Shape `mapstructure:"shapes" json:"shapes"`

	// Position is where this object is located in game units, where 1
	// game unit is 64 pixels at standard zoom.
	Position Vector `mapstructure:"position" json:"position"`

	// Velocity is the change in position of the object in game units per
	// second.
	Velocity Vector `mapstructure:"velocity" json:"velocity"`

	// Rotation is the rotation of the object in radians.
	Rotation float64 `mapstructure:"rotation" json:"rotation"`

	// AngularVelocity is the change in rotation of the object in
	// radians per second
	AngularVelocity float64 `mapstructure:"angular_velocity" json:"angular_velocity"`
}
