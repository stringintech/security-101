package com.stringintech.security101;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.stringintech.security101.dto.UserLoginRequest;
import com.stringintech.security101.dto.UserRegisterRequest;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
class DtoValidationIntegrationTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Test
    void whenRegisterWithEmptyFields_thenReturn400AndValidationErrors() throws Exception {
        UserRegisterRequest request = new UserRegisterRequest();
        // All fields empty

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.message").value("Validation failed"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='fullName')].message").value("Full name is required"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='username')].message").value("Username is required"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='password')].message").value("Password is required"));
    }

    @Test
    void whenRegisterWithInvalidFields_thenReturn400AndValidationErrors() throws Exception {
        UserRegisterRequest request = new UserRegisterRequest();
        request.setFullName("A"); // Too short
        request.setUsername("a"); // Too short, invalid format
        request.setPassword("weak"); // Too short and missing required characters

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.message").value("Validation failed"))
                .andExpect(jsonPath("$.timestamp").exists())
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='fullName')].message")
                        .value("Full name must be between 2 and 100 characters"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='username')].message")
                        .value("Username must be 3-20 characters long and contain only letters, numbers, underscores and hyphens"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='password')].message")
                        .value(containsInAnyOrder(
                                "Password must be at least 8 characters long",
                                "Password must contain at least one digit, one lowercase, one uppercase, and one special character"
                        )));
    }

    @Test
    void whenRegisterWithValidFields_thenReturn200() throws Exception {
        UserRegisterRequest request = new UserRegisterRequest();
        request.setFullName("Test User");
        request.setUsername("testuser");
        request.setPassword("Test123@pass");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk());
    }

    @Test
    void whenLoginWithEmptyFields_thenReturn400AndValidationErrors() throws Exception {
        UserLoginRequest request = new UserLoginRequest();
        // All fields empty

        mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.message").value("Validation failed"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='username')].message").value("Username is required"))
                .andExpect(jsonPath("$.fieldErrors[?(@.field=='password')].message").value("Password is required"));
    }

    @Test
    void whenLoginWithValidFields_thenReturn200() throws Exception {
        // First register a user
        UserRegisterRequest registerRequest = new UserRegisterRequest();
        registerRequest.setFullName("Test User");
        registerRequest.setUsername("testlogin");
        registerRequest.setPassword("Test123@pass");

        mockMvc.perform(post("/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(registerRequest)))
                .andExpect(status().isOk());

        // Then try to login
        UserLoginRequest loginRequest = new UserLoginRequest();
        loginRequest.setUsername("testlogin");
        loginRequest.setPassword("Test123@pass");

        mockMvc.perform(post("/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isOk());
    }
}