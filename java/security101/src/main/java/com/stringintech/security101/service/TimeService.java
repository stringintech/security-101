package com.stringintech.security101.service;

import org.springframework.stereotype.Service;

import java.time.Clock;
import java.time.Instant;

@Service
public class TimeService {

    private final Clock clock;

    public TimeService(Clock clock) {
        this.clock = clock;
    }

    public Instant now() {
        return clock.instant();
    }
}
