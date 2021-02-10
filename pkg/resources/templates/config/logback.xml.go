// Copyright 2020 Orange SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package apis

package config

var LogbackTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!--
  Licensed to the Apache Software Foundation (ASF) under one or more
  contributor license agreements.  See the NOTICE file distributed with
  this work for additional information regarding copyright ownership.
  The ASF licenses this file to You under the Apache License, Version 2.0
  (the "License"); you may not use this file except in compliance with
  the License.  You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
-->
<configuration scan="true" scanPeriod="30 seconds">
    <contextListener class="ch.qos.logback.classic.jul.LevelChangePropagator">
        <resetJUL>true</resetJUL>
    </contextListener>
    
    <appender name="APP_FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>${org.apache.nifi.bootstrap.config.log.dir}/nifi-app.log</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.SizeAndTimeBasedRollingPolicy">
            <!--
              For daily rollover, use 'app_%d.log'.
              For hourly rollover, use 'app_%d{yyyy-MM-dd_HH}.log'.
              To GZIP rolled files, replace '.log' with '.log.gz'.
              To ZIP rolled files, replace '.log' with '.log.zip'.
            -->
            <fileNamePattern>${org.apache.nifi.bootstrap.config.log.dir}/nifi-app_%d{yyyy-MM-dd_HH}.%i.log</fileNamePattern>
            <maxFileSize>{{.MaxFileSizeAppFile}}</maxFileSize>
            <!-- keep 30 log files worth of history -->
            <maxHistory>{{.MaxHistoryAppFile}}</maxHistory>
        </rollingPolicy>
        <immediateFlush>true</immediateFlush>
        <encoder class="ch.qos.logback.classic.encoder.PatternLayoutEncoder">
            <pattern>%date %level [%thread] %logger{40} %msg%n</pattern>
        </encoder>
    </appender>
    
    <appender name="USER_FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>${org.apache.nifi.bootstrap.config.log.dir}/nifi-user.log</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <!--
              For daily rollover, use 'user_%d.log'.
              For hourly rollover, use 'user_%d{yyyy-MM-dd_HH}.log'.
              To GZIP rolled files, replace '.log' with '.log.gz'.
              To ZIP rolled files, replace '.log' with '.log.zip'.
            -->
            <fileNamePattern>${org.apache.nifi.bootstrap.config.log.dir}/nifi-user_%d.log</fileNamePattern>
            <!-- keep 30 log files worth of history -->
            <maxHistory>{{.MaxHistoryUserFile}}</maxHistory>
        </rollingPolicy>
        <encoder class="ch.qos.logback.classic.encoder.PatternLayoutEncoder">
            <pattern>%date %level [%thread] %logger{40} %msg%n</pattern>
        </encoder>
    </appender>

    <appender name="BOOTSTRAP_FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>${org.apache.nifi.bootstrap.config.log.dir}/nifi-bootstrap.log</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <!--
              For daily rollover, use 'user_%d.log'.
              For hourly rollover, use 'user_%d{yyyy-MM-dd_HH}.log'.
              To GZIP rolled files, replace '.log' with '.log.gz'.
              To ZIP rolled files, replace '.log' with '.log.zip'.
            -->
            <fileNamePattern>${org.apache.nifi.bootstrap.config.log.dir}/nifi-bootstrap_%d.log</fileNamePattern>
            <!-- keep 5 log files worth of history -->
            <maxHistory>{{.MaxHistoryBootstrapFile}}</maxHistory>
        </rollingPolicy>
        <encoder class="ch.qos.logback.classic.encoder.PatternLayoutEncoder">
            <pattern>%date %level [%thread] %logger{40} %msg%n</pattern>
        </encoder>
    </appender>
	
    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder class="ch.qos.logback.classic.encoder.PatternLayoutEncoder">
            <pattern>%date %level [%thread] %logger{40} %msg%n</pattern>
        </encoder>
    </appender>
    
    <!-- valid logging levels: TRACE, DEBUG, INFO, WARN, ERROR -->
    
    <logger name="org.apache.nifi" level="{{.LogLevelNifi}}"/>
    <logger name="org.apache.nifi.processors" level="{{.LogLevelNifiProcessors}}"/>
    <logger name="org.apache.nifi.processors.standard.LogAttribute" level="{{.LogLevelNifiProcessorsStandardLogAttribute}}"/>
    <logger name="org.apache.nifi.processors.standard.LogMessage" level="{{.LogLevelNifiProcessorsStandardLogMessage}}"/>
    <logger name="org.apache.nifi.controller.repository.StandardProcessSession" level="{{.LogLevelNifiControllerRepositoryStandardProcessSession}}" />
    
    
    <logger name="org.apache.zookeeper.ClientCnxn" level="{{.LogLevelZookeeperClientCnxn}}" />
    <logger name="org.apache.zookeeper.server.NIOServerCnxn" level="{{.LogLevelZookeeperServerNIOServerCnxn}}" />
    <logger name="org.apache.zookeeper.server.NIOServerCnxnFactory" level="{{.LogLevelZookeeperServerNIOServerCnxnFactory}}" />
    <logger name="org.apache.zookeeper.server.quorum" level="{{.LogLevelZookeeperServerQuorum}}" />
    <logger name="org.apache.zookeeper.ZooKeeper" level="{{.LogLevelZookeeperZooKeeper}}" />
    <logger name="org.apache.zookeeper.server.PrepRequestProcessor" level="{{.LogLevelZookeeperServerPrepRequestProcessor}}" />

    <logger name="org.apache.calcite.runtime.CalciteException" level="{{.LogLevelCalciteRuntimeCalciteException}}" />

    <logger name="org.apache.curator.framework.recipes.leader.LeaderSelector" level="{{.LogLevelCuratorFrameworkRecipesLeaderLeaderSelector}}" />
    <logger name="org.apache.curator.ConnectionState" level="{{.LogLevelCuratorConnectionState}}" />
    
    <!-- Logger for managing logging statements for nifi clusters. -->
    <logger name="org.apache.nifi.cluster" level="{{.LogLevelNifiCluster}}"/>

    <!-- Logger for logging HTTP requests received by the web server. -->
    <logger name="org.apache.nifi.server.JettyServer" level="{{.LogLevelNifiServerJettyServer}}"/>

    <!-- Logger for managing logging statements for jetty -->
    <logger name="org.eclipse.jetty" level="{{.LogLevelJetty}}"/>

    <!-- Suppress non-error messages due to excessive logging by class or library -->
    <logger name="org.springframework" level="{{.LogLevelSpringframework}}"/>
    
    <!-- Suppress non-error messages due to known warning about redundant path annotation (NIFI-574) -->
    <logger name="org.glassfish.jersey.internal.Errors" level="{{.LogLevelJerseyInternalErrors}}"/>

    <!--
        Logger for capturing user events. We do not want to propagate these
        log events to the root logger. These messages are only sent to the
        user-log appender.
    -->
    <logger name="org.apache.nifi.web.security" level="{{.LogLevelNifiWebSecurity}}" additivity="false">
        <appender-ref ref="USER_FILE"/>
    </logger>
    <logger name="org.apache.nifi.web.api.config" level="{{.LogLevelNifiWebApiConfig}}" additivity="false">
        <appender-ref ref="USER_FILE"/>
    </logger>
    <logger name="org.apache.nifi.authorization" level="{{.LogLevelNifiAuthorization}}" additivity="false">
        <appender-ref ref="USER_FILE"/>
    </logger>
    <logger name="org.apache.nifi.cluster.authorization" level="{{.LogLevelNifiClusterAuthorization}}" additivity="false">
        <appender-ref ref="USER_FILE"/>
    </logger>
    <logger name="org.apache.nifi.web.filter.RequestLogger" level="{{.LogLevelNifiWebFilterRequestLogger}}" additivity="false">
        <appender-ref ref="USER_FILE"/>
    </logger>


    <!--
        Logger for capturing Bootstrap logs and NiFi's standard error and standard out. 
    -->
    <logger name="org.apache.nifi.bootstrap" level="{{.LogLevelNifiBootstrap}}" additivity="false">
        <appender-ref ref="BOOTSTRAP_FILE" />
    </logger>
    <logger name="org.apache.nifi.bootstrap.Command" level="{{.LogLevelNifiBootstrapCommand}}" additivity="false">
        <appender-ref ref="CONSOLE" />
        <appender-ref ref="BOOTSTRAP_FILE" />
    </logger>

    <!-- Everything written to NiFi's Standard Out will be logged with the logger org.apache.nifi.StdOut at INFO level -->
    <logger name="org.apache.nifi.StdOut" level="{{.LogLevelNifiStdOut}}" additivity="false">
        <appender-ref ref="BOOTSTRAP_FILE" />
    </logger>
    
    <!-- Everything written to NiFi's Standard Error will be logged with the logger org.apache.nifi.StdErr at ERROR level -->
    <logger name="org.apache.nifi.StdErr" level="{{.LogLevelNifiStdErr}}" additivity="false">
        <appender-ref ref="BOOTSTRAP_FILE" />
    </logger>


    <root level="{{.LogLevelRoot}}">
        <appender-ref ref="APP_FILE"/>
    </root>
    
</configuration>
`
