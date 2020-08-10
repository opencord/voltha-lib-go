/* 
 * Copyright 2018-present Open Networking Foundation
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package Config provides dynamic log correlation enable/disable for specific Voltha component with logcorrelation lookup
// from etcd kvstore implemented using Backend.
// Any Voltha component can start utilizing dynamic log correlation feature by starting goroutine of StartLogCorrelationConfigProcessing after
// starting kvClient for the component.

package config

import (
	"context"
	"errors"
	"github.com/opencord/voltha-lib-go/v3/pkg/log"
	"os"
)

// Flag indicating the current log correlation status
var logCorrelationStatus bool = true

const defaultLogCorrelationKey = "default"   // kvstore key containing default log correlation status

// StartLogCorrelationConfigProcessiong initialize compoment config
type ComponentLogCorrelationController struct {
        ComponentName         string
        componentNameConfig   *ComponentConfig
        configManager         *ConfigManager
        initialLogCorrelation  string // Initial default log correlation
}

func NewComponentLogCorrelationController(ctx context.Context, cm *ConfigManager) (*ComponentLogCorrelationController, error) {
        logger.Debug(ctx, "creating-new-component-log-correlation-controller")
        componentName := os.Getenv("COMPONENT_NAME")
        if componentName == "" {
                return nil, errors.New("Unable to retrieve PoD Component Name from Runtime env")
        }

        return &ComponentLogCorrelationController{
                ComponentName:         componentName,
                componentNameConfig:   nil,
                configManager:         cm,
                initialLogCorrelation: "ENABLED",
        }, nil
}

func StartLogCorrelationConfigProcessing(cm *ConfigManager, ctx context.Context) {
        cc, err := NewComponentLogCorrelationController(ctx, cm)
        if err != nil {
                logger.Errorw(ctx, "unnable-to-construct-component-log-controller-instance-for-log-correlation-config-monitoring", log.Fields{"error": err})
                return
        }

        cc.componentNameConfig = cm.InitComponentConfig(cc.ComponentName, ConfigTypeLogCorrelation)
        logger.Debugw(ctx, "component-log-correlation-config", log.Fields{"cc-component-name-config": cc.componentNameConfig})

	cc.persistDefaultLogCorrelationStatus(ctx)
        cc.processLogCorrelationConfig(ctx)
}

func (c *ComponentLogCorrelationController) persistDefaultLogCorrelationStatus(ctx context.Context) {
	_, err := c.componentNameConfig.Retrieve(ctx, defaultLogCorrelationKey)
	if err != nil {
		logger.Debugw(ctx, "failed-to-retrieve-component-default-log-correlation-status-at-startup", log.Fields{"error": err})

		err = c.componentNameConfig.Save(ctx, defaultLogCorrelationKey, c.initialLogCorrelation)
		if err != nil {
			logger.Errorw(ctx, "failed-to-persist-component-default-log-correlation-status-at-startup", log.Fields{"error": err, "logcorrelation": c.initialLogCorrelation})
		}
	}
}

func (c *ComponentLogCorrelationController) processLogCorrelationConfig(ctx context.Context) {

	// Load and apply Log Correlation Status for first time
	initialLogCorrelationStatus, err := c.getInitialLogCorrelationStatus(ctx)
	if err != nil {
		logger.Warnw(ctx, "unable-to-load-log-correlation-config-at-startup", log.Fields{"error": err})
	} else {
		if err := c.updateLogCorrelationStatus(ctx, initialLogCorrelationStatus); err != nil {
			logger.Warnw(ctx, "unable-to-apply-log-config-at-startup", log.Fields{"error": err})
		}
	}

        componentConfigEventChan := c.componentNameConfig.MonitorForConfigChange(ctx)

        // process the events for componentName
        var configEvent *ConfigChangeEvent
        for {
                configEvent = <-componentConfigEventChan
                logger.Debugw(ctx, "processing-log-config-change", log.Fields{"ChangeType": configEvent.ChangeType, "Package": configEvent.ConfigAttribute})

                logCorrelationConfig, err := c.getLogCorrelationStatus(ctx)
                if err != nil {
                        logger.Warnw(ctx, "unable-to-get-log-correlation-status", log.Fields{"error": err})
                        continue
                }

                if err := c.updateLogCorrelationStatus(ctx, logCorrelationConfig); err != nil {
                        logger.Warnw(ctx, "unable-to-update-log-correlation-status", log.Fields{"error": err})
                }
        }
}

func (c *ComponentLogCorrelationController) getInitialLogCorrelationStatus(ctx context.Context) (map[string]string, error) {
        logger.Debug(ctx, "get-initial-log-correlation-status")

        initialLogCorrelationStatus, err := c.componentNameConfig.Retrieve(ctx, defaultLogCorrelationKey)
        if err != nil {
                return nil, err
        }

	logConfig := make(map[string]string)
	logConfig[c.ComponentName] = initialLogCorrelationStatus

	return logConfig, err
}

func (c *ComponentLogCorrelationController) getLogCorrelationStatus(ctx context.Context) (map[string]string, error) {
        logger.Debug(ctx, "get-log-correlation-status")

        componentLogCorrelationConfig, err := c.componentNameConfig.RetrieveAll(ctx)
	if err != nil {
                return nil, err
        }
        return componentLogCorrelationConfig, err
}

func (c *ComponentLogCorrelationController) updateLogCorrelationStatus(ctx context.Context, componentLogCorrelationConfig map[string]string) error {
        logger.Debug(ctx, "update-log-correlation-status")

	currentLogCorrelationStatus := c.getStatus(ctx, componentLogCorrelationConfig)

        if logCorrelationStatus != currentLogCorrelationStatus {
                log.SetLogCorrelation(ctx, currentLogCorrelationStatus)
                logCorrelationStatus = currentLogCorrelationStatus
        } else {
                logger.Debug(ctx, "logcorrelation-status-same-as-currently-active")
        }
        return nil
}

func (c *ComponentLogCorrelationController) getStatus(ctx context.Context, logCorrelationConfig map[string]string) bool {
        logger.Debug(ctx, "get-status")

        for _, status := range logCorrelationConfig {
                logger.Debugw(ctx, "log-correlation-status", log.Fields{"status": status})
                if status == "ENABLED" {
                        logger.Debugw(ctx, "log-correlation-status-enabled", log.Fields{"status": status})
                        return true
                } else if status == "DISABLED" {
                        logger.Debugw(ctx, "log-correlation-status-disabled", log.Fields{"status": status})
                        return false
                }
        }
        return true
}

