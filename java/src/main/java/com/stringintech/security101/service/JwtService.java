package com.stringintech.security101.service;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Decoders;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKey;
import java.time.Instant;
import java.util.Date;
import java.util.concurrent.TimeUnit;

@Service
public class JwtService {

    private final String secretKey;
    private final Integer expirationDays;
    private final TimeService timeService;

    public JwtService(
            @Value("${security.jwt.secret-key}") String secretKey,
            @Value("${security.jwt.expiration-days}") Integer expirationDays, TimeService timeService
    ) {
        this.secretKey = secretKey;
        this.expirationDays = expirationDays;
        this.timeService = timeService;
    }

    public String extractUsernameFromAndValidate(String token) {
        Claims claims;
        try {
            claims = extractAllClaims(token);
            if (claims == null) {
                throw new JwtException("Failed to parse token");
            }
        } catch (Exception e) {
            throw new JwtException("Failed to parse token", e);
        }
        String username = claims.getSubject();
        if (username == null) {
            throw new JwtException("Failed to extract username from token");
        }
        Date creationTime = claims.getIssuedAt();
        if (creationTime == null) {
            throw new JwtException("Failed to extract creation time from token");
        }
        if (isTokenExpired(creationTime)) {
            throw new JwtException("Token has expired");
        }
        return username;
    }

    public String generateToken(UserDetails userDetails, Instant creationTime) {
        return Jwts
                .builder()
                .subject(userDetails.getUsername())
                .issuedAt(Date.from(creationTime))
                .signWith(getSignInKey())
                .compact();
    }

    private Claims extractAllClaims(String token) {
        //TODO parse might throw exception
        Object payload = Jwts.parser().verifyWith(getSignInKey()).build().parse(token).getPayload();
        return payload instanceof Claims ? (Claims) payload : null; //TODO when null?
    }

    private SecretKey getSignInKey() {
        byte[] keyBytes = Decoders.BASE64.decode(secretKey);
        return Keys.hmacShaKeyFor(keyBytes);
    }

    private boolean isTokenExpired(Date creationTime) {
        Date expirationDate = new Date(creationTime.getTime() + TimeUnit.DAYS.toMillis(expirationDays));
        return expirationDate.before(Date.from(timeService.now()));
    }
}
