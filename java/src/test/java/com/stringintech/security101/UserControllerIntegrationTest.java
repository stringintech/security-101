package com.stringintech.security101;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.stringintech.security101.dto.UserLoginRequest;
import com.stringintech.security101.dto.UserRegisterRequest;
import com.stringintech.security101.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.boot.test.mock.mockito.SpyBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import java.time.Clock;
import java.time.Instant;
import java.time.ZoneId;

import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
class UserControllerIntegrationTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @MockBean
    private Clock clock;

    @SpyBean
    private UserRepository userRepository;

    private static final Instant FIXED_TIME = Instant.parse("2024-01-01T10:00:00Z");

    @BeforeEach
    void setUp() {
        when(clock.instant()).thenReturn(FIXED_TIME);
        when(clock.getZone()).thenReturn(ZoneId.systemDefault());
    }

    @Test
        // Request hits SecurityFilterChain first, which allows /auth/** paths without authentication
        // Request goes to UserController.register(), which encodes password using BCrypt
        // Creates and stores new user via UserRepository
        // Returns UserDto without sensitive password field
    void whenRegisterUser_thenReturn200AndUser() throws Exception {
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("testuser");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Test User");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.username").value("testuser"))
                .andExpect(jsonPath("$.fullName").value("Test User"))
                .andExpect(jsonPath("$.password").doesNotExist());
    }

    @Test
        // Same flow as above, but UserRepository.createUser() throws DuplicateUsernameException
        // GlobalExceptionHandler catches this and returns 400 with error details
    void whenRegisterExistingUsername_thenReturn400() throws Exception {
        // First registration
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("duplicate");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Test User");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk());

        // Duplicate registration
        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isBadRequest());
    }

    @Test
        // First register call succeeds as above
        // Login request hits SecurityFilterChain (allowed for /auth/**)
        // UserController.login() uses AuthenticationManager to verify credentials
        // On success, JwtService generates a token using the fixed clock time
        // Returns token + user details
    void whenLoginValidUser_thenReturnTokenAndUser() throws Exception {
        // First register a user
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("logintest");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Login Test");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk());

        // Then try to login
        UserLoginRequest loginDto = new UserLoginRequest();
        loginDto.setUsername("logintest");
        loginDto.setPassword("Test123@pass");

        mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginDto)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.token").exists())
                .andExpect(jsonPath("$.user.username").value("logintest"))
                .andExpect(jsonPath("$.tokenCreationTime").value(FIXED_TIME.toString()));
    }

    @Test
        // AuthenticationManager fails to authenticate since user doesn't exist
        // Spring's DaoAuthenticationProvider throws AuthenticationException
        // Custom JwtAuthenticationEntryPoint returns 401
    void whenLoginWithNonexistentUsername_thenReturn401() throws Exception {
        // Then try to login with wrong password
        UserLoginRequest loginDto = new UserLoginRequest();
        loginDto.setUsername("nonexistent");
        loginDto.setPassword("password");

        mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginDto)))
                .andExpect(status().isUnauthorized());
    }

    @Test
        // Same flow as above, but AuthenticationManager fails password comparison
        // using BCryptPasswordEncoder, resulting in 401
    void whenLoginWithWrongPassword_thenReturn401() throws Exception {
        // First register a user
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("wrongpass");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Wrong Pass Test");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk());

        // Then try to login with wrong password
        UserLoginRequest loginDto = new UserLoginRequest();
        loginDto.setUsername("wrongpass");
        loginDto.setPassword("wrongpassword");

        mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginDto)))
                .andExpect(status().isUnauthorized());
    }

    @Test
        // Register & login succeed as above
        // /users/me request hits JwtAuthenticationFilter first
        // Filter validates token and loads user via UserDetailsService
        // Sets SecurityContext authentication
        // Request proceeds to controller which returns current user
    void whenGetUserWithValidToken_thenReturnUser() throws Exception {
        // First register and login to get token
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("getuser");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Get User Test");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk());

        UserLoginRequest loginDto = new UserLoginRequest();
        loginDto.setUsername("getuser");
        loginDto.setPassword("Test123@pass");

        String response = mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginDto)))
                .andExpect(status().isOk())
                .andReturn()
                .getResponse()
                .getContentAsString();

        String token = objectMapper.readTree(response).get("token").asText();

        // Then try to get user details
        mockMvc.perform(post("/users/me")
                        .header("Authorization", "Bearer " + token))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.username").value("getuser"));
    }

    @Test
        // JwtAuthenticationFilter checks Authorization header
        // Missing/invalid header causes filter to proceed to next filter in chain
        // Which results in 401 from custom JwtAuthenticationEntryPoint
    void whenGetUserWithoutAuthHeader_thenReturn401() throws Exception {
        mockMvc.perform(post("/users/me"))
                .andExpect(status().isUnauthorized());
    }

    @Test
        // JwtAuthenticationFilter attempts to parse invalid token
        // JwtService throws JwtException which is caught by filter
        // Filter silently continues to next filter without setting SecurityContext
        // Leading to 401 from custom JwtAuthenticationEntryPoint
    void whenGetUserWithInvalidTokenFormat_thenReturn401() throws Exception {
        mockMvc.perform(post("/users/me")
                        .header("Authorization", "Bearer invalid.token.format"))
                .andExpect(status().isUnauthorized());
    }

    @Test
        // Token appears valid but JwtService.isTokenExpired() returns true
        // Because mock clock is set 11 days ahead
        // JwtService throws expired token exception which is caught by filter
        // Filter silently continues to next filter without setting SecurityContext
        // Leading to 401 from custom JwtAuthenticationEntryPoint
    void whenGetUserWithExpiredToken_thenReturn401() throws Exception {
        // First register and login to get token
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("expiredtoken");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Expired Token Test");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk());

        UserLoginRequest loginDto = new UserLoginRequest();
        loginDto.setUsername("expiredtoken");
        loginDto.setPassword("Test123@pass");

        String response = mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginDto)))
                .andExpect(status().isOk())
                .andReturn()
                .getResponse()
                .getContentAsString();

        String token = objectMapper.readTree(response).get("token").asText();

        // Move time forward past token expiration
        when(clock.instant()).thenReturn(FIXED_TIME.plusSeconds(60 * 60 * 24 * 11)); // 11 days later

        // Try to get user details with expired token
        mockMvc.perform(post("/users/me")
                        .header("Authorization", "Bearer " + token))
                .andExpect(status().isUnauthorized());
    }

    @Test
        // Token format and expiry are valid, but mocked UserRepository returns null
        // This causes UserDetailsService to throw UsernameNotFoundException
        // Filter silently continues to next filter without setting SecurityContext
        // Leading to 401 from custom JwtAuthenticationEntryPoint
    void whenGetUserWithNonExistentUsername_thenReturn401() throws Exception {
        // First register and login to get token
        UserRegisterRequest registerDto = new UserRegisterRequest();
        registerDto.setUsername("deleteuser");
        registerDto.setPassword("Test123@pass");
        registerDto.setFullName("Delete User Test");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerDto)))
                .andExpect(status().isOk());

        UserLoginRequest loginDto = new UserLoginRequest();
        loginDto.setUsername("deleteuser");
        loginDto.setPassword("Test123@pass");

        String response = mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginDto)))
                .andExpect(status().isOk())
                .andReturn()
                .getResponse()
                .getContentAsString();

        String token = objectMapper.readTree(response).get("token").asText();

        // Only now override the repository behavior to simulate deleted user
        doReturn(null).when(userRepository).getUserByUsername("deleteuser");

        // Try to get user details with token that has valid format but non-existent user
        mockMvc.perform(post("/users/me")
                        .header("Authorization", "Bearer " + token))
                .andExpect(status().isUnauthorized());

        // Verify the repository was called
        verify(userRepository, atLeastOnce()).getUserByUsername("deleteuser");

        // Reset the spy to not affect other tests
        reset(userRepository);
    }
}