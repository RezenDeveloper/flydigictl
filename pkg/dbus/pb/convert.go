package pb

import (
	"strconv"

	"github.com/pipe01/flydigictl/pkg/flydigi/config"
	"golang.org/x/image/colornames"
)

func ColorFromRGB(r, g, b byte) *Color {
	return &Color{Rgb: int32(b) | (int32(g) << 8) | (int32(r) << 16)}
}

func ColorFromHex(hex string) (*Color, bool) {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return nil, false
	}

	n, err := strconv.ParseUint(hex, 16, 32)
	if err != nil || n > 0xFFFFFF {
		return nil, false
	}

	return &Color{Rgb: int32(n)}, true
}

func ColorFromName(name string) (*Color, bool) {
	rgba, ok := colornames.Map[name]
	if !ok {
		return nil, false
	}

	return ColorFromRGB(rgba.R, rgba.G, rgba.B), true
}

func (c *Color) RGB() (r, g, b byte) {
	return byte(c.Rgb >> 16), byte(c.Rgb >> 8), byte(c.Rgb)
}

func (c *Color) LedUnit() *config.LedUnit {
	r, g, b := c.RGB()
	return &config.LedUnit{r, g, b}
}

func ConvertGamepadConfiguration(bean *config.AllConfigBean) *GamepadConfiguration {
	return &GamepadConfiguration{
		LeftJoystick:  ConvertJoystickConfiguration(bean.JoyMapping.LeftJoystic),
		RightJoystick: ConvertJoystickConfiguration(bean.JoyMapping.RightJoystic),
		LeftTrigger:   ConvertTriggerConfiguration(bean.TriggerMapping.LeftTrigger),
		RightTrigger:  ConvertTriggerConfiguration(bean.TriggerMapping.RightTrigger),
	}
}

func (c *GamepadConfiguration) ApplyTo(bean *config.AllConfigBean) {
	c.LeftJoystick.ApplyTo(bean.JoyMapping.LeftJoystic)
	c.RightJoystick.ApplyTo(bean.JoyMapping.RightJoystic)
	c.LeftTrigger.ApplyTo(bean.TriggerMapping.LeftTrigger)
	c.RightTrigger.ApplyTo(bean.TriggerMapping.RightTrigger)
}

func ConvertJoystickConfiguration(bean *config.JoyStickBean) *JoystickConfiguration {
	return &JoystickConfiguration{
		Deadzone: bean.Curve.Zero,
	}
}

func (c *JoystickConfiguration) ApplyTo(bean *config.JoyStickBean) {
	bean.Curve.Zero = c.Deadzone
}

func ConvertTriggerConfiguration(bean *config.Trigger) *TriggerConfiguration {
	return &TriggerConfiguration{
		AutoTrigger: &AutoTrigger{
			Mode: bean.AutoTrigger.Mode,
			VibrationBind: &VibrationBind{
				Type:          bean.AutoTrigger.VibrationBind.Type,
				MinFilter:     bean.AutoTrigger.VibrationBind.MinFilter,
				Scale:         bean.AutoTrigger.VibrationBind.Scale,
				TriggerParams: bean.AutoTrigger.VibrationBind.TriggerParams,
			},
			MixedBorder: bean.AutoTrigger.MixedBorder,
			MixedParams: bean.AutoTrigger.MixedParams,
		},
		TriggerMotor: &TriggerMotor{
			LineGear: &TriggerMotorSet{
				Type:      bean.TriggerMotor.LineGear.Type,
				Min:       bean.TriggerMotor.LineGear.Min,
				Max:       bean.TriggerMotor.LineGear.Max,
				Filter:    bean.TriggerMotor.LineGear.Filter,
				VibrLimit: bean.TriggerMotor.LineGear.Vibrlimit,
				Scale:     bean.TriggerMotor.LineGear.Scale,
				TimeLimit: bean.TriggerMotor.LineGear.TimeLimit,
			},
			MicrGear: &TriggerMotorSet{
				Type:      bean.TriggerMotor.MicrGear.Type,
				Min:       bean.TriggerMotor.MicrGear.Min,
				Max:       bean.TriggerMotor.MicrGear.Max,
				Filter:    bean.TriggerMotor.MicrGear.Filter,
				VibrLimit: bean.TriggerMotor.MicrGear.Vibrlimit,
				Scale:     bean.TriggerMotor.MicrGear.Scale,
				TimeLimit: bean.TriggerMotor.MicrGear.TimeLimit,
			},
		},
	}
}

func (c *TriggerConfiguration) ApplyTo(bean *config.Trigger) {
	if bean == nil {
		return
	}

	bean.Type = 0
	bean.AutoTrigger = &config.Autotrigger{
		Mode: c.AutoTrigger.Mode,
		VibrationBind: &config.Vibrationbind{
			Type:          c.AutoTrigger.VibrationBind.Type,
			MinFilter:     c.AutoTrigger.VibrationBind.MinFilter,
			Scale:         c.AutoTrigger.VibrationBind.Scale,
			TriggerParams: c.AutoTrigger.VibrationBind.TriggerParams,
		},
		MixedBorder: c.AutoTrigger.MixedBorder,
		MixedParams: c.AutoTrigger.MixedParams,
	}
	bean.TriggerMotor = &config.TriggerMotor{
		LineGear: &config.TriggerMotorSet{
			Type:      c.TriggerMotor.LineGear.Type,
			Min:       c.TriggerMotor.LineGear.Min,
			Max:       c.TriggerMotor.LineGear.Max,
			Filter:    c.TriggerMotor.LineGear.Filter,
			Vibrlimit: c.TriggerMotor.LineGear.VibrLimit,
			Scale:     c.TriggerMotor.LineGear.Scale,
			TimeLimit: c.TriggerMotor.LineGear.TimeLimit,
		},
		MicrGear: &config.TriggerMotorSet{
			Type:      c.TriggerMotor.MicrGear.Type,
			Min:       c.TriggerMotor.MicrGear.Min,
			Max:       c.TriggerMotor.MicrGear.Max,
			Filter:    c.TriggerMotor.MicrGear.Filter,
			Vibrlimit: c.TriggerMotor.MicrGear.VibrLimit,
			Scale:     c.TriggerMotor.MicrGear.Scale,
			TimeLimit: c.TriggerMotor.MicrGear.TimeLimit,
		},
	}
}

func ConvertLEDConfiguration(bean *config.NewLedConfigBean) *LedsConfiguration {
	var leds isLedsConfiguration_Leds

	switch bean.LedMode {
	case config.LedModeOff:
		leds = &LedsConfiguration_Off{}

	case config.LedModeSteady:
		unit := bean.LedGroups[0].Units[0]

		leds = &LedsConfiguration_Steady{
			Steady: &LedsSteady{
				Color: ColorFromRGB(unit.R, unit.G, unit.B),
			},
		}

	case config.LedModeStreamlined:
		leds = &LedsConfiguration_Streamlined{
			Streamlined: &LedsStreamlined{
				Speed: float32(100-bean.Loop_time) / 100,
			},
		}
	case config.LedModeGradient:
		colors := make([]*Color, 0)

		if bean.Rgb_num > 0 && len(bean.LedGroups) > 0 {
			group := bean.LedGroups[0]

			for i := 0; i < int(bean.Rgb_num) && i < len(group.Units); i++ {
				unit := group.Units[i]
				if unit == nil {
					continue
				}

				colors = append(colors, &Color{
					Rgb: int32(unit.R)<<16 | int32(unit.G)<<8 | int32(unit.B),
				})
			}
		}

		leds = &LedsConfiguration_Gradient{
			Gradient: &LedsGradient{
				Speed:  float32(100-bean.Loop_time) / 100,
				Colors: colors,
			},
		}

	default:
		panic("TODO: implement")
	}

	return &LedsConfiguration{
		Leds:       leds,
		Brightness: float32(bean.Light_scale) / 255,
	}
}

func (c *LedsConfiguration) ApplyTo(bean *config.NewLedConfigBean) {
	bean.Light_scale = byte(c.Brightness * 100)

	switch leds := c.Leds.(type) {
	case *LedsConfiguration_Off:

	case *LedsConfiguration_Steady:
		bean.SetSteady(*leds.Steady.Color.LedUnit())

	case *LedsConfiguration_Streamlined:
		bean.SetStreamlined(leds.Streamlined.Speed)

	case *LedsConfiguration_Gradient:
		g := leds.Gradient

		colors := make([]config.LedUnit, 0, len(g.Colors))
		for _, c := range g.Colors {
			if c == nil {
				continue
			}
			colors = append(colors, *c.LedUnit())
		}

		if len(colors) < 2 {
			return
		}

		bean.SetGradient(colors, g.Speed)
	}
}
