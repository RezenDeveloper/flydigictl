package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"github.com/spf13/cobra"
)

type triggerSide string

type TriggerMotorSetJson *struct {
	Type      *int32 `json:"Type"`
	Min       *int32 `json:"Min"`
	Max       *int32 `json:"Max"`
	Filter    *int32 `json:"Filter"`
	VibrLimit *int32 `json:"VibrLimit"`
	Scale     *int32 `json:"Scale"`
	TimeLimit *int32 `json:"TimeLimit"`
}

type TriggerJSON struct {
	AutoTrigger *struct {
		Mode *int32 `json:"Mode"`

		VibrationBind *struct {
			Type          *int32  `json:"Type"`
			MinFilter     *int32  `json:"MinFilter"`
			Scale         *int32  `json:"Scale"`
			TriggerParams []int32 `json:"TriggerParams"`
		} `json:"VibrationBind"`

		MixedBorder *int32  `json:"MixedBorder"`
		MixedParams []int32 `json:"MixedParams"`
	} `json:"AutoTrigger"`

	TriggerMotor *struct {
		LineGear TriggerMotorSetJson `json:"LineGear"`
		MicrGear TriggerMotorSetJson `json:"MicrGear"`
	} `json:"TriggerMotor"`
}

const (
	triggerLeft  triggerSide = "left"
	triggerRight triggerSide = "right"
)

func (s triggerSide) GetBean(b *pb.GamepadConfiguration) *pb.TriggerConfiguration {
	switch s {
	case triggerLeft:
		return b.LeftTrigger
	case triggerRight:
		return b.RightTrigger
	}
	panic("invalid trigger side")
}

func triggerModeName(mode int32) string {
	switch mode {
	case 0:
		return "default"
	case 1:
		return "race"
	case 2:
		return "recoil"
	case 3:
		return "sniper"
	case 4:
		return "lock"
	case 5:
		return "vibration"
	default:
		return fmt.Sprintf("unknown(%d)", mode)
	}
}

func loadTriggerJSON(path string) (*TriggerJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg TriggerJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyTriggerJSON(dst *pb.TriggerConfiguration, src *TriggerJSON) {
	if src.AutoTrigger != nil {
		at := dst.AutoTrigger

		if src.AutoTrigger.Mode != nil {
			at.Mode = *src.AutoTrigger.Mode
		}

		if src.AutoTrigger.VibrationBind != nil {
			vb := at.VibrationBind

			if src.AutoTrigger.VibrationBind.Type != nil {
				vb.Type = *src.AutoTrigger.VibrationBind.Type
			}
			if src.AutoTrigger.VibrationBind.MinFilter != nil {
				vb.MinFilter = *src.AutoTrigger.VibrationBind.MinFilter
			}
			if src.AutoTrigger.VibrationBind.Scale != nil {
				vb.Scale = *src.AutoTrigger.VibrationBind.Scale
			}
			if len(src.AutoTrigger.VibrationBind.TriggerParams) > 0 {
				vb.TriggerParams = src.AutoTrigger.VibrationBind.TriggerParams
			}
		}

		if src.AutoTrigger.MixedBorder != nil {
			at.MixedBorder = *src.AutoTrigger.MixedBorder
		}
		if len(src.AutoTrigger.MixedParams) > 0 {
			at.MixedParams = src.AutoTrigger.MixedParams
		}
	}

	if src.TriggerMotor != nil {
		if src.TriggerMotor.LineGear != nil {
			lg := dst.TriggerMotor.LineGear

			if src.TriggerMotor.LineGear.Type != nil {
				lg.Type = *src.TriggerMotor.LineGear.Type
			}
			if src.TriggerMotor.LineGear.Min != nil {
				lg.Min = *src.TriggerMotor.LineGear.Min
			}
			if src.TriggerMotor.LineGear.Max != nil {
				lg.Max = *src.TriggerMotor.LineGear.Max
			}
			if src.TriggerMotor.LineGear.Filter != nil {
				lg.Filter = *src.TriggerMotor.LineGear.Filter
			}
			if src.TriggerMotor.LineGear.VibrLimit != nil {
				lg.VibrLimit = *src.TriggerMotor.LineGear.VibrLimit
			}
			if src.TriggerMotor.LineGear.Scale != nil {
				lg.Scale = *src.TriggerMotor.LineGear.Scale
			}
			if src.TriggerMotor.LineGear.TimeLimit != nil {
				lg.TimeLimit = *src.TriggerMotor.LineGear.TimeLimit
			}
		}
		if src.TriggerMotor.MicrGear != nil {
			lg := dst.TriggerMotor.MicrGear

			if src.TriggerMotor.MicrGear.Type != nil {
				lg.Type = *src.TriggerMotor.MicrGear.Type
			}
			if src.TriggerMotor.MicrGear.Min != nil {
				lg.Min = *src.TriggerMotor.MicrGear.Min
			}
			if src.TriggerMotor.MicrGear.Max != nil {
				lg.Max = *src.TriggerMotor.MicrGear.Max
			}
			if src.TriggerMotor.MicrGear.Filter != nil {
				lg.Filter = *src.TriggerMotor.MicrGear.Filter
			}
			if src.TriggerMotor.MicrGear.VibrLimit != nil {
				lg.VibrLimit = *src.TriggerMotor.MicrGear.VibrLimit
			}
			if src.TriggerMotor.MicrGear.Scale != nil {
				lg.Scale = *src.TriggerMotor.MicrGear.Scale
			}
			if src.TriggerMotor.MicrGear.TimeLimit != nil {
				lg.TimeLimit = *src.TriggerMotor.MicrGear.TimeLimit
			}
		}
	}
}

var triggerCommand = &cobra.Command{
	Use:   "trigger",
	Short: "Configure triggers (Apex 4)",
}

func genTriggerCommand(side triggerSide) *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(side),
		Short: fmt.Sprintf("Manage %s trigger configuration", side),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := readConfiguration(func(conf *pb.GamepadConfiguration) int32 {
				return int32(side.GetBean(conf).AutoTrigger.Mode)
			})
			if err != nil {
				return fmt.Errorf("get configuration: %w", err)
			}

			modeName := triggerModeName(mode)

			if terseOutput {
				fmt.Println(modeName)
			} else {
				fmt.Printf("%s trigger mode: %s\n", side, modeName)
			}
			return nil
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "default",
		Short: "Set this trigger to default mode",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.AutoTrigger.Mode = 0
				cfg.AutoTrigger.VibrationBind.TriggerParams = []int32{1, 10, 1, 90, 0}

				for i := range cfg.AutoTrigger.MixedParams {
					cfg.AutoTrigger.MixedParams[i] = 0
				}

				lg := cfg.TriggerMotor.LineGear
				lg.Type = 1
				lg.Min = 30
				lg.Max = 80
				lg.Filter = 5
				lg.VibrLimit = 1
				lg.Scale = 50
				lg.TimeLimit = 0
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "race",
		Short: "Set this trigger to race mode",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.AutoTrigger.Mode = 1
				cfg.AutoTrigger.VibrationBind.TriggerParams = []int32{100, 1, 255, 70, 0}

				for i := range cfg.AutoTrigger.MixedParams {
					cfg.AutoTrigger.MixedParams[i] = 0
				}
				cfg.AutoTrigger.MixedParams[1] = 30

				lg := cfg.TriggerMotor.LineGear
				lg.Type = 1
				lg.Min = 90
				lg.Max = 100
				lg.Filter = 5
				lg.VibrLimit = 1
				lg.Scale = 100
				lg.TimeLimit = 0
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "recoil",
		Short: "Set this trigger to recoil mode",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.AutoTrigger.Mode = 2
				cfg.AutoTrigger.VibrationBind.Type = 0
				cfg.AutoTrigger.VibrationBind.MinFilter = 10
				cfg.AutoTrigger.VibrationBind.Scale = 50
				cfg.AutoTrigger.VibrationBind.TriggerParams =
					[]int32{100, 1, 255, 70, 0}

				cfg.AutoTrigger.MixedBorder = 0
				cfg.AutoTrigger.MixedParams = []int32{
					0, 1, 50, 15, 1,
					0, 0, 0, 0, 0,
				}

				lg := cfg.TriggerMotor.LineGear
				lg.Type = 1
				lg.Min = 30
				lg.Max = 80
				lg.Filter = 5
				lg.VibrLimit = 1
				lg.Scale = 50
				lg.TimeLimit = 0
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "sniper",
		Short: "Set this trigger to sniper mode",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.AutoTrigger.Mode = 3
				cfg.AutoTrigger.VibrationBind.Type = 0
				cfg.AutoTrigger.VibrationBind.MinFilter = 10
				cfg.AutoTrigger.VibrationBind.Scale = 50
				cfg.AutoTrigger.VibrationBind.TriggerParams =
					[]int32{100, 1, 255, 70, 0}

				cfg.AutoTrigger.MixedBorder = 0
				cfg.AutoTrigger.MixedParams = []int32{
					50, 30, 1, 0, 1,
					0, 0, 0, 0, 0,
				}

				lg := cfg.TriggerMotor.LineGear
				lg.Type = 1
				lg.Min = 30
				lg.Max = 80
				lg.Filter = 5
				lg.VibrLimit = 1
				lg.Scale = 50
				lg.TimeLimit = 0
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "lock",
		Short: "Set this trigger to lock mode",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.AutoTrigger.Mode = 4
				cfg.AutoTrigger.VibrationBind.Type = 0
				cfg.AutoTrigger.VibrationBind.MinFilter = 10
				cfg.AutoTrigger.VibrationBind.Scale = 50
				cfg.AutoTrigger.VibrationBind.TriggerParams =
					[]int32{100, 1, 255, 70, 0}

				cfg.AutoTrigger.MixedBorder = 0
				cfg.AutoTrigger.MixedParams = []int32{
					40, 250, 1, 0, 0,
					0, 0, 0, 0, 0,
				}

				lg := cfg.TriggerMotor.LineGear
				lg.Type = 1
				lg.Min = 30
				lg.Max = 80
				lg.Filter = 5
				lg.VibrLimit = 1
				lg.Scale = 50
				lg.TimeLimit = 0
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "vibration",
		Short: "Set this trigger to vibration mode",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.AutoTrigger.Mode = 5
				cfg.AutoTrigger.VibrationBind.Type = 2
				cfg.AutoTrigger.VibrationBind.MinFilter = 10
				cfg.AutoTrigger.VibrationBind.Scale = 50
				cfg.AutoTrigger.VibrationBind.TriggerParams =
					[]int32{1, 1, 1, 90, 0}

				cfg.AutoTrigger.MixedBorder = 0
				cfg.AutoTrigger.MixedParams = []int32{
					1, 1, 1, 90, 0,
					0, 0, 0, 0, 0,
				}

				lg := cfg.TriggerMotor.LineGear
				lg.Type = 1
				lg.Min = 30
				lg.Max = 80
				lg.Filter = 5
				lg.VibrLimit = 1
				lg.Scale = 50
				lg.TimeLimit = 0
			})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "custom <path>",
		Short: "Load custom trigger configuration from JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			info, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("file does not exist: %s", path)
				}
				return fmt.Errorf("cannot access file %s: %w", path, err)
			}

			if info.IsDir() {
				return fmt.Errorf("path is a directory, not a file: %s", path)
			}

			if filepath.Ext(path) != ".json" {
				return fmt.Errorf("expected a .json file: %s", path)
			}

			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				jsonCfg, err := loadTriggerJSON(path)
				if err != nil {
					panic(err)
				}
				applyTriggerJSON(cfg, jsonCfg)
			})
		},
	})

	return cmd
}

func init() {
	var triggerLeftCommand = genTriggerCommand(triggerLeft)
	var triggerRightCommand = genTriggerCommand(triggerRight)

	triggerCommand.AddCommand(triggerLeftCommand)
	triggerCommand.AddCommand(triggerRightCommand)

	rootCmd.AddCommand(triggerCommand)
}
