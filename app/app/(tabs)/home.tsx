import React, { useEffect, useState } from "react";
import {
  Text,
  View,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
  ScrollView,
  Alert,
} from "react-native";
import { getUserInfo, logout } from "../../lib/api";
import { Card, Icon, Button } from "@rneui/themed";
import { router, useLocalSearchParams } from "expo-router";

interface UserInfoResponse {
  status: number;
  data: {
    attributes: {
      email: string;
      email_verified: string;
      name: string;
      sub: string;
    };
    username: string;
  };
}

export default function Home() {
  const [userInfo, setUserInfo] = useState<UserInfoResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);
  const [loggingOut, setLoggingOut] = useState(false);

  const params = useLocalSearchParams();
  const isNewLogin = params.newLogin === "true";

  const fetchUserInfo = async () => {
    try {
      setLoading(true);
      setError(null);

      if (isNewLogin) {
        await new Promise((resolve) => setTimeout(resolve, 500));
      }

      const data = await getUserInfo();
      console.log("User data retrieved successfully:", data.status);
      setUserInfo(data);
    } catch (err: any) {
      console.error("Error fetching user info:", err);

      if (err.response) {
        console.error("Response data:", err.response.data);
        console.error("Response status:", err.response.status);
      }

      if (err.response?.status === 401) {
        setError("Your session has expired. Please log in again.");
      } else {
        setError(err.message || "Failed to fetch user information");
      }
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await fetchUserInfo();
  };

  useEffect(() => {
    fetchUserInfo();
  }, []);

  const handleLogout = async () => {
    Alert.alert("Logout", "Are you sure you want to logout?", [
      {
        text: "Cancel",
        style: "cancel",
      },
      {
        text: "Yes, Logout",
        style: "destructive",
        onPress: async () => {
          try {
            setLoggingOut(true);
            await logout();
            router.replace("/(auth)");
          } catch (error) {
            console.error("Logout error:", error);
            Alert.alert("Error", "Failed to logout. Please try again.");
          } finally {
            setLoggingOut(false);
          }
        },
      },
    ]);
  };

  const renderUserInfo = () => {
    if (loading) {
      return (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#3b82f6" />
          <Text style={styles.loadingText}>Loading user information...</Text>
        </View>
      );
    }

    if (error) {
      return (
        <View style={styles.errorContainer}>
          <Icon
            name="error-outline"
            type="material"
            size={48}
            color="#ef4444"
          />
          <Text style={styles.errorText}>{error}</Text>
          <Button
            title="Try Again"
            onPress={fetchUserInfo}
            buttonStyle={styles.retryButton}
          />
        </View>
      );
    }

    if (!userInfo) {
      return (
        <View style={styles.errorContainer}>
          <Text style={styles.noDataText}>No user information available</Text>
        </View>
      );
    }

    const { attributes, username } = userInfo.data;

    return (
      <Card containerStyle={styles.card}>
        <Card.Title style={styles.cardTitle}>User Profile</Card.Title>
        <Card.Divider />

        <View style={styles.profileHeader}>
          <Icon
            name="account-circle"
            type="material"
            size={80}
            color="#3b82f6"
          />
          <Text style={styles.nameText}>{attributes.name}</Text>
        </View>

        <View style={styles.infoSection}>
          <Card.Divider />

          <View style={styles.infoRow}>
            <Icon name="email" type="material" size={20} color="#6b7280" />
            <Text style={styles.infoLabel}>Email:</Text>
            <Text style={styles.infoValue}>{attributes.email}</Text>
          </View>

          <View style={styles.infoRow}>
            <Icon
              name={
                attributes.email_verified === "true" ? "check-circle" : "cancel"
              }
              type="material"
              size={20}
              color={
                attributes.email_verified === "true" ? "#22c55e" : "#ef4444"
              }
            />
            <Text style={styles.infoLabel}>Email Verified:</Text>
            <Text
              style={[
                styles.infoValue,
                {
                  color:
                    attributes.email_verified === "true"
                      ? "#22c55e"
                      : "#ef4444",
                },
              ]}
            >
              {attributes.email_verified === "true" ? "Yes" : "No"}
            </Text>
          </View>

          <View style={styles.infoRow}>
            <Icon
              name="fingerprint"
              type="material"
              size={20}
              color="#6b7280"
            />
            <Text style={styles.infoLabel}>User ID:</Text>
            <Text style={styles.infoValue}>{attributes.sub}</Text>
          </View>
        </View>
      </Card>
    );
  };

  return (
    <ScrollView
      contentContainerStyle={styles.container}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={onRefresh}
          colors={["#3b82f6"]}
        />
      }
    >
      {renderUserInfo()}

      <Button
        title={loggingOut ? "Logging out..." : "Logout"}
        icon={
          <Icon
            name="logout"
            type="material"
            size={20}
            color="white"
            style={{ marginRight: 8 }}
          />
        }
        buttonStyle={styles.logoutButton}
        onPress={handleLogout}
        disabled={loggingOut}
      />

      <Text style={styles.refreshHint}>
        Pull down to refresh user information
      </Text>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flexGrow: 1,
    padding: 16,
    backgroundColor: "#f9fafb",
  },
  loadingContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  loadingText: {
    marginTop: 12,
    fontSize: 16,
    color: "#6b7280",
  },
  errorContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 20,
  },
  errorText: {
    marginTop: 12,
    fontSize: 16,
    color: "#ef4444",
    textAlign: "center",
  },
  noDataText: {
    fontSize: 16,
    color: "#6b7280",
  },
  retryButton: {
    marginTop: 16,
    backgroundColor: "#3b82f6",
    paddingHorizontal: 20,
  },
  card: {
    borderRadius: 12,
    padding: 16,
    marginVertical: 16,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  cardTitle: {
    fontSize: 22,
    color: "#111827",
  },
  profileHeader: {
    alignItems: "center",
    marginVertical: 16,
  },
  nameText: {
    fontSize: 24,
    fontWeight: "bold",
    color: "#111827",
    marginTop: 12,
  },
  infoSection: {
    marginTop: 8,
  },
  infoRow: {
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: 12,
  },
  infoLabel: {
    fontSize: 16,
    color: "#4b5563",
    fontWeight: "500",
    marginLeft: 8,
    width: 110,
  },
  infoValue: {
    fontSize: 16,
    color: "#111827",
    flex: 1,
  },
  logoutButton: {
    backgroundColor: "#ef4444",
    borderRadius: 8,
    marginVertical: 16,
    height: 48,
  },
  refreshHint: {
    textAlign: "center",
    color: "#6b7280",
    fontSize: 14,
    marginVertical: 16,
  },
});
